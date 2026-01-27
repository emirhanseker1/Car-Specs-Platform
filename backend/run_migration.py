import sqlite3
import sys

def run_migration():
    """Run the 4-level migration script on the database"""
    
    # Read the migration SQL
    with open('backend/migrations/006_migrate_to_4_level.sql', 'r', encoding='utf-8') as f:
        migration_sql = f.read()
    
    # Connect to database
    conn = sqlite3.connect('backend/vehicles.db')
    cursor = conn.cursor()
    
    print("=" * 70)
    print("STARTING MIGRATION: 3-Level → 4-Level Structure")
    print("=" * 70)
    print()
    
    try:
        # Execute the migration script
        cursor.executescript(migration_sql)
        conn.commit()
        
        print("✓ Migration executed successfully!")
        print()
        
        # Verification queries
        print("=" * 70)
        print("VERIFICATION RESULTS")
        print("=" * 70)
        print()
        
        # Count records at each level
        print("📊 Record Counts:")
        print("-" * 70)
        cursor.execute("""
            SELECT 'Brands' as level, COUNT(*) as count FROM brands
            UNION ALL
            SELECT 'Models', COUNT(*) FROM models
            UNION ALL
            SELECT 'Generations', COUNT(*) FROM generations
            UNION ALL
            SELECT 'Trims', COUNT(*) FROM trims
        """)
        for row in cursor.fetchall():
            print(f"  {row[0]:<15} {row[1]:>5}")
        print()
        
        # Check for orphaned trims
        cursor.execute("SELECT COUNT(*) FROM trims WHERE generation_id IS NULL")
        orphaned = cursor.fetchone()[0]
        if orphaned == 0:
            print(f"✓ All trims linked to generations (0 orphaned)")
        else:
            print(f"⚠ WARNING: {orphaned} trims without generation_id!")
        print()
        
        # Show sample hierarchy
        print("=" * 70)
        print("SAMPLE HIERARCHY (First 10 Generations)")
        print("=" * 70)
        print()
        cursor.execute("""
            SELECT 
                b.name as brand,
                m.name as model,
                g.code as gen_code,
                g.name as generation,
                g.start_year || '-' || COALESCE(CAST(g.end_year AS TEXT), 'Present') as years,
                COUNT(t.id) as trims
            FROM brands b
            JOIN models m ON m.brand_id = b.id
            JOIN generations g ON g.model_id = m.id
            LEFT JOIN trims t ON t.generation_id = g.id
            GROUP BY b.id, m.id, g.id
            ORDER BY b.name, m.name, g.start_year
            LIMIT 10
        """)
        
        print(f"{'Brand':<12} {'Model':<15} {'Code':<8} {'Generation':<20} {'Years':<15} {'Trims':>6}")
        print("-" * 80)
        for row in cursor.fetchall():
            print(f"{row[0]:<12} {row[1]:<15} {row[2]:<8} {row[3] or '':<20} {row[4]:<15} {row[5]:>6}")
        
        print()
        print("=" * 70)
        print("✓ MIGRATION COMPLETED SUCCESSFULLY!")
        print("=" * 70)
        
    except Exception as e:
        conn.rollback()
        print(f"✗ Migration failed: {e}")
        sys.exit(1)
    finally:
        conn.close()

if __name__ == "__main__":
    run_migration()
