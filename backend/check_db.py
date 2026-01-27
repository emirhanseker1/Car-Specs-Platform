import sqlite3

conn = sqlite3.connect('backend/vehicles.db')
cursor = conn.cursor()

print("=" * 70)
print("DATABASE STATUS CHECK")
print("=" * 70)
print()

# Count records
cursor.execute('SELECT COUNT(*) FROM brands')
print(f"Brands: {cursor.fetchone()[0]}")

cursor.execute('SELECT COUNT(*) FROM models')
print(f"Models: {cursor.fetchone()[0]}")

cursor.execute('SELECT COUNT(*) FROM generations')
print(f"Generations: {cursor.fetchone()[0]}")

cursor.execute('SELECT COUNT(*) FROM trims')
print(f"Trims: {cursor.fetchone()[0]}")

print()
cursor.execute('SELECT COUNT(*) FROM trims WHERE generation_id IS NULL')
orphaned = cursor.fetchone()[0]
if orphaned == 0:
    print(f"✓ All trims linked to generations (0 orphaned)")
else:
    print(f"⚠ WARNING: {orphaned} trims without generation_id!")

print()
print("Sample data:")
print("-" * 70)
cursor.execute('''
    SELECT 
        b.name as brand,
        m.name as model,
        g.code as gen_code,
        t.name as trim_name,
        t.year
    FROM trims t
    LEFT JOIN generations g ON t.generation_id = g.id
    LEFT JOIN models m ON g.model_id = m.id
    LEFT JOIN brands b ON m.brand_id = b.id
    LIMIT 5
''')

for row in cursor.fetchall():
    print(f"{row[0]:10} | {row[1]:15} | {row[2]:8} | {row[3]:20} | {row[4]}")

conn.close()
