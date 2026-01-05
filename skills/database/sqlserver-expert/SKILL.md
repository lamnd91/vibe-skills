---
name: sqlserver-expert
description: Expert in Microsoft SQL Server development and administration. Use when writing T-SQL queries, stored procedures, optimizing database performance (deadlocks, slow queries, execution plans), designing schemas, configuring SQL Server, implementing CDC (Change Data Capture), or integrating SQL Server with .NET Core/C# using Entity Framework Core or Dapper.
---

# SQL Server Expert

Act as DBA and developer expert in Microsoft SQL Server.

## T-SQL Advanced Patterns

**CTEs with Window Functions:**
```sql
WITH RankedData AS (
  SELECT Id, Name, Department,
    ROW_NUMBER() OVER (PARTITION BY Department ORDER BY HireDate) AS RowNum,
    SUM(Salary) OVER (PARTITION BY Department) AS DeptTotal,
    LAG(Salary) OVER (ORDER BY HireDate) AS PrevSalary
  FROM Employees
)
SELECT * FROM RankedData WHERE RowNum = 1;
```

**MERGE Statement:**
```sql
MERGE INTO Target AS t
USING Source AS s ON t.Id = s.Id
WHEN MATCHED THEN UPDATE SET t.Name = s.Name, t.ModifiedAt = GETDATE()
WHEN NOT MATCHED THEN INSERT (Id, Name) VALUES (s.Id, s.Name)
WHEN NOT MATCHED BY SOURCE THEN DELETE;
```

**Pagination:**
```sql
SELECT * FROM Orders
ORDER BY OrderDate DESC
OFFSET @PageSize * (@PageNumber - 1) ROWS
FETCH NEXT @PageSize ROWS ONLY;
```

**APPLY Operators:**
```sql
-- Get top 3 orders per customer
SELECT c.Name, o.OrderDate, o.TotalAmount
FROM Customers c
CROSS APPLY (
  SELECT TOP 3 * FROM Orders WHERE CustomerId = c.Id ORDER BY OrderDate DESC
) o;
```

## .NET Core Integration

### Entity Framework Core Setup
```csharp
// DbContext
public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options) : base(options) { }
    
    public DbSet<Order> Orders => Set<Order>();
    public DbSet<Customer> Customers => Set<Customer>();
    
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Order>(e =>
        {
            e.ToTable("Orders");
            e.HasKey(x => x.Id);
            e.Property(x => x.TotalAmount).HasColumnType("decimal(18,2)");
            e.HasIndex(x => x.CustomerId);
        });
    }
}

// Registration in Program.cs
builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("Default"),
        sqlOptions =>
        {
            sqlOptions.EnableRetryOnFailure(maxRetryCount: 3);
            sqlOptions.CommandTimeout(30);
        }));
```

### Raw SQL & Stored Procedures (EF Core)
```csharp
// Raw SQL query
var orders = await context.Orders
    .FromSqlInterpolated($"SELECT * FROM Orders WHERE CustomerId = {customerId}")
    .ToListAsync();

// Stored procedure
var results = await context.Database
    .SqlQueryRaw<OrderSummary>("EXEC dbo.GetOrderSummary @CustomerId = {0}", customerId)
    .ToListAsync();

// Execute non-query
await context.Database.ExecuteSqlInterpolatedAsync(
    $"UPDATE Orders SET Status = {newStatus} WHERE Id = {orderId}");
```

### Dapper Integration
```csharp
// Setup
public class DapperContext
{
    private readonly string _connectionString;
    public DapperContext(IConfiguration config) 
        => _connectionString = config.GetConnectionString("Default")!;
    
    public IDbConnection CreateConnection() => new SqlConnection(_connectionString);
}

// Query
public async Task<IEnumerable<Order>> GetOrdersAsync(int customerId)
{
    using var conn = _context.CreateConnection();
    return await conn.QueryAsync<Order>(
        "SELECT * FROM Orders WHERE CustomerId = @CustomerId",
        new { CustomerId = customerId });
}

// Stored procedure
public async Task<OrderSummary> GetSummaryAsync(int customerId)
{
    using var conn = _context.CreateConnection();
    return await conn.QueryFirstOrDefaultAsync<OrderSummary>(
        "dbo.GetOrderSummary",
        new { CustomerId = customerId },
        commandType: CommandType.StoredProcedure);
}

// Bulk insert with transaction
public async Task BulkInsertAsync(IEnumerable<Order> orders)
{
    using var conn = _context.CreateConnection();
    conn.Open();
    using var tx = conn.BeginTransaction();
    try
    {
        await conn.ExecuteAsync(
            @"INSERT INTO Orders (CustomerId, OrderDate, TotalAmount) 
              VALUES (@CustomerId, @OrderDate, @TotalAmount)",
            orders, transaction: tx);
        tx.Commit();
    }
    catch { tx.Rollback(); throw; }
}
```

### SqlBulkCopy for Large Datasets
```csharp
public async Task BulkCopyAsync(DataTable data, string tableName)
{
    using var conn = new SqlConnection(_connectionString);
    await conn.OpenAsync();
    using var bulk = new SqlBulkCopy(conn)
    {
        DestinationTableName = tableName,
        BatchSize = 10000,
        BulkCopyTimeout = 600
    };
    await bulk.WriteToServerAsync(data);
}
```

## Performance Tuning

### Execution Plan Analysis
```sql
-- Top CPU-consuming queries
SELECT TOP 10
  total_worker_time/execution_count AS avg_cpu_time,
  execution_count,
  SUBSTRING(st.text, (qs.statement_start_offset/2)+1, 100) AS query_text
FROM sys.dm_exec_query_stats qs
CROSS APPLY sys.dm_exec_sql_text(qs.sql_handle) st
ORDER BY total_worker_time DESC;
```

### Missing Indexes
```sql
SELECT 
  migs.avg_total_user_cost * migs.avg_user_impact * (migs.user_seeks + migs.user_scans) AS improvement,
  mid.statement AS table_name,
  mid.equality_columns,
  mid.inequality_columns,
  mid.included_columns
FROM sys.dm_db_missing_index_groups mig
JOIN sys.dm_db_missing_index_group_stats migs ON mig.index_group_handle = migs.group_handle
JOIN sys.dm_db_missing_index_details mid ON mig.index_handle = mid.index_handle
ORDER BY improvement DESC;
```

### Index Fragmentation
```sql
SELECT 
  OBJECT_NAME(ips.object_id) AS table_name,
  i.name AS index_name,
  ips.avg_fragmentation_in_percent,
  CASE WHEN ips.avg_fragmentation_in_percent > 30 THEN 'REBUILD'
       WHEN ips.avg_fragmentation_in_percent > 10 THEN 'REORGANIZE'
       ELSE 'OK' END AS action
FROM sys.dm_db_index_physical_stats(DB_ID(), NULL, NULL, NULL, 'LIMITED') ips
JOIN sys.indexes i ON ips.object_id = i.object_id AND ips.index_id = i.index_id
WHERE ips.avg_fragmentation_in_percent > 10 AND ips.page_count > 1000;
```

### Currently Running & Blocking
```sql
SELECT 
  r.session_id, r.status, r.wait_type, r.wait_time, r.blocking_session_id,
  SUBSTRING(st.text, (r.statement_start_offset/2)+1, 200) AS query_text
FROM sys.dm_exec_requests r
CROSS APPLY sys.dm_exec_sql_text(r.sql_handle) st
WHERE r.session_id > 50
ORDER BY r.total_elapsed_time DESC;
```

## Deadlock Detection & Resolution

### Extended Events for Deadlocks
```sql
CREATE EVENT SESSION [DeadlockCapture] ON SERVER
ADD EVENT sqlserver.xml_deadlock_report
ADD TARGET package0.event_file(SET filename=N'C:\Temp\Deadlocks.xel')
WITH (MAX_MEMORY=4096 KB);
ALTER EVENT SESSION [DeadlockCapture] ON SERVER STATE = START;
```

### Query Deadlock History
```sql
SELECT 
  xed.value('@timestamp', 'datetime2(3)') AS deadlock_time,
  xed.query('.') AS deadlock_graph
FROM (
  SELECT CAST(target_data AS XML) AS target_data
  FROM sys.dm_xe_session_targets st
  JOIN sys.dm_xe_sessions s ON s.address = st.event_session_address
  WHERE s.name = 'system_health' AND st.target_name = 'ring_buffer'
) AS data
CROSS APPLY target_data.nodes('RingBufferTarget/event[@name="xml_deadlock_report"]') AS xed(xed)
ORDER BY deadlock_time DESC;
```

### Prevention Strategies
1. Access tables in consistent order across all procedures
2. Keep transactions short
3. Use READ COMMITTED SNAPSHOT isolation
4. Add proper indexes to reduce lock duration
5. Use `NOLOCK` or `READPAST` hints where appropriate

## Change Data Capture (CDC)

### Setup CDC
```sql
-- Enable on database
EXEC sys.sp_cdc_enable_db;

-- Enable on table
EXEC sys.sp_cdc_enable_table
  @source_schema = N'dbo',
  @source_name = N'Orders',
  @role_name = NULL,
  @supports_net_changes = 1,
  @captured_column_list = N'Id,CustomerId,OrderDate,TotalAmount,Status';
```

### Query CDC Changes
```sql
DECLARE @from_lsn binary(10) = sys.fn_cdc_get_min_lsn('dbo_Orders');
DECLARE @to_lsn binary(10) = sys.fn_cdc_get_max_lsn();

SELECT 
  CASE __$operation WHEN 1 THEN 'DELETE' WHEN 2 THEN 'INSERT' 
       WHEN 3 THEN 'BEFORE_UPDATE' WHEN 4 THEN 'AFTER_UPDATE' END AS op,
  sys.fn_cdc_map_lsn_to_time(__$start_lsn) AS change_time,
  Id, CustomerId, TotalAmount
FROM cdc.fn_cdc_get_all_changes_dbo_Orders(@from_lsn, @to_lsn, N'all');
```

### CDC Consumer in C#
```csharp
public class CdcService
{
    private readonly string _connString;
    private byte[]? _lastLsn;
    
    public async Task<IEnumerable<CdcChange>> GetChangesAsync(string captureInstance)
    {
        using var conn = new SqlConnection(_connString);
        
        // Get LSN range
        var minLsn = _lastLsn ?? await conn.ExecuteScalarAsync<byte[]>(
            $"SELECT sys.fn_cdc_get_min_lsn('{captureInstance}')");
        var maxLsn = await conn.ExecuteScalarAsync<byte[]>(
            "SELECT sys.fn_cdc_get_max_lsn()");
        
        if (minLsn.SequenceEqual(maxLsn)) return [];
        
        var changes = await conn.QueryAsync<CdcChange>(
            $@"SELECT __$operation AS Operation, 
                      sys.fn_cdc_map_lsn_to_time(__$start_lsn) AS ChangeTime, *
               FROM cdc.fn_cdc_get_all_changes_{captureInstance}(@from, @to, N'all')",
            new { from = minLsn, to = maxLsn });
        
        _lastLsn = maxLsn;
        return changes;
    }
}
```

## System Queries

### Table Structure
```sql
SELECT c.COLUMN_NAME, c.DATA_TYPE, c.CHARACTER_MAXIMUM_LENGTH, c.IS_NULLABLE
FROM INFORMATION_SCHEMA.COLUMNS c
WHERE c.TABLE_SCHEMA = @schema AND c.TABLE_NAME = @table
ORDER BY c.ORDINAL_POSITION;
```

### All Indexes on Table
```sql
SELECT i.name, i.type_desc, i.is_unique, i.is_primary_key,
  STRING_AGG(c.name, ', ') WITHIN GROUP (ORDER BY ic.key_ordinal) AS columns
FROM sys.indexes i
JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
JOIN sys.columns c ON ic.object_id = c.object_id AND ic.column_id = c.column_id
WHERE i.object_id = OBJECT_ID(@tableName)
GROUP BY i.name, i.type_desc, i.is_unique, i.is_primary_key;
```

### Foreign Keys
```sql
SELECT fk.name AS fk_name, tp.name AS parent_table, cp.name AS parent_column,
       tr.name AS ref_table, cr.name AS ref_column
FROM sys.foreign_keys fk
JOIN sys.foreign_key_columns fkc ON fk.object_id = fkc.constraint_object_id
JOIN sys.tables tp ON fkc.parent_object_id = tp.object_id
JOIN sys.columns cp ON fkc.parent_object_id = cp.object_id AND fkc.parent_column_id = cp.column_id
JOIN sys.tables tr ON fkc.referenced_object_id = tr.object_id
JOIN sys.columns cr ON fkc.referenced_object_id = cr.object_id AND fkc.referenced_column_id = cr.column_id
WHERE tp.name = @tableName;
```

### Table Sizes
```sql
SELECT s.name + '.' + t.name AS table_name, p.rows,
  CAST(SUM(a.total_pages) * 8 / 1024.0 AS DECIMAL(10,2)) AS total_mb
FROM sys.tables t
JOIN sys.schemas s ON t.schema_id = s.schema_id
JOIN sys.indexes i ON t.object_id = i.object_id
JOIN sys.partitions p ON i.object_id = p.object_id AND i.index_id = p.index_id
JOIN sys.allocation_units a ON p.partition_id = a.container_id
WHERE i.index_id <= 1
GROUP BY s.name, t.name, p.rows
ORDER BY SUM(a.total_pages) DESC;
```

## Best Practices

### Performance
1. Avoid `SELECT *` - list columns explicitly
2. Use appropriate indexes for WHERE/JOIN columns
3. Avoid functions on columns in WHERE (not sargable)
4. Use `SET NOCOUNT ON` in stored procedures
5. Use `OPTION (RECOMPILE)` for parameter-sensitive queries

### Security
1. Never concatenate strings - use parameters
2. Least privilege for application users
3. Use schemas to organize and control access