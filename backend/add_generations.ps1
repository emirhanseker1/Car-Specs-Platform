$dbPath = "car_specs.db"

# SQL to insert generations
$sql = @"
INSERT INTO generations (model_id, code, name, start_year, end_year, is_current) VALUES
(1, '8L', 'Typ 8L', 1996, 2003, 0),
(1, '8P', 'Typ 8P', 2003, 2012, 0),
(1, '8V', 'Typ 8V', 2012, 2020, 0),
(1, '8Y', 'Typ 8Y', 2020, NULL, 1);
"@

# Try different approaches to execute SQL
Write-Host "🚀 Adding Audi A3 Generations..." -ForegroundColor Cyan

# Approach 1: Using .NET SQLite (if available)
try {
    Add-Type -Path "System.Data.SQLite.dll" -ErrorAction SilentlyContinue
    $connection = New-Object -TypeName System.Data.SQLite.SQLiteConnection
    $connection.ConnectionString = "Data Source=$dbPath"
    $connection.Open()
    
    $command = $connection.CreateCommand()
    $command.CommandText = $sql
    $command.ExecuteNonQuery()
    
    $connection.Close()
    Write-Host "✅ Successfully added generations!" -ForegroundColor Green
    
    # Verify
    $connection.Open()
    $verifyCmd = $connection.CreateCommand()
    $verifyCmd.CommandText = "SELECT COUNT(*) FROM generations WHERE model_id = 1"
    $count = $verifyCmd.ExecuteScalar()
    $connection.Close()
    Write-Host "📊 Total generations for Audi A3: $count" -ForegroundColor Yellow
    exit 0
}
catch {
    Write-Host "⚠️ .NET SQLite not available, trying alternative method..." -ForegroundColor Yellow
}

# Approach 2: Direct file write and manual instruction
Write-Host ""
Write-Host "📝 SQL generated. Please run manually:" -ForegroundColor Yellow
Write-Host $sql
Write-Host ""
Write-Host "Or use a SQLite client to execute migrations/008_add_audi_a3_generations.sql" -ForegroundColor Yellow
