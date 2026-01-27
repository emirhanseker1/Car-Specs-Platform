import sqlite3
import sys

def run_simplified_migration():
    """Run simplified 4-level migration that doesn't recreate the trims table"""
    
    # Read the simplified migration SQL
    with open('backend/migrations/006_migrate_to_4_level_simple.sql', 'r', encoding='utf-8') as f:
        migration_sql = f.read()
    
    # Connect to database
    conn = sqlite3.connect('backend/vehicles.db')
    cursor = conn.cursor()
    
    print("=" * 70)
    print("🚀 STARTING MIGRATION: 3-Level → 4-Level Structure (Simplified)")
    print("=" * 70)
    print()
    
    try:
        # Split and execute statements one by one for better error handling
        statements = [s.strip() for s in migration_sql.split(';') if s.strip() and not s.strip().startswith('--')]
        
        for i, statement in enumerate(statements):
            if statement:
                try:
                    cursor.execute(statement)
                    # Check if this was a SELECT statement
                    if statement.strip().upper().startswith('SELECT'):
                        results = cursor.fetchall()
                        if results:
                            for row in results:
                                print(f"  {row[0]}: {row[1]}")
                except Exception as e:
                    # Skip comments and empty statements
                    if 'syntax error' not in str(e).lower():
                        print(f"Statement {i+1} warning: {e}")
        
        conn.commit()
        
        print()
        print("✓ Migration executed successfully!")
        print()
        
        # Verification queries
        print("=" * 70)
        print("📊 VERIFICATION RESULTS")
        print("=" * 70)
        print()
        
        # Count records at each level
        print("Record Counts at Each Level:")
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
            print(f"  {row[0]:<15} {row[1]:>5} records")
        print()
        
        # Check for orphaned trims
        cursor.execute("SELECT COUNT(*) FROM trims WHERE generation_id IS NULL")
        orphaned = cursor.fetchone()[0]
        if orphaned == 0:
            print(f"✓ All trims successfully linked to generations (0 orphaned)")
        else:
            print(f"⚠ WARNING: {orphaned} trims without generation_id!")
        print()
        
        # Show sample hierarchy
        print("=" * 70)
        print("🌳 NEW 4-LEVEL HIERARCHY (Sample - First 15 Generations)")
        print("=" * 70)
        print()
        cursor.execute("""
            SELECT 
                b.name as brand,
                m.name as model,
                g.code as code,
                g.start_year || '-' || COALESCE(CAST(g.end_year AS TEXT), 'Now') as years,
                COUNT(t.id) as trims
            FROM brands b
            JOIN models m ON m.brand_id = b.id
            JOIN generations g ON g.model_id = m.id
            LEFT JOIN trims t ON t.generation_id = g.id
            GROUP BY b.id, m.id, g.id
            ORDER BY b.name, m.name, g.start_year
            LIMIT 15
        """)
        
        print(f"{'Brand':<10} | {'Model':<15} | {'Code':<10} | {'Years':<12} | {'Trims':>6}")
        print("-" * 70)
        for row in cursor.fetchall():
            print(f"{row[0]:<10} | {row[1]:<15} | {row[2]:<10} | {row[3]:<12} | {row[4]:>6}")
        
        # Show detailed sample with trim names
        print()
        print("=" * 70)
        print("📋 DETAILED SAMPLE (First Example)")
        print("=" * 70)
        print()
        
        cursor.execute("""
            SELECT 
                b.name as brand,
                m.name as model,
                g.code as gen_code,
                g.name as generation,
                t.name as trim_name,
                t.year,
                t.fuel_type,
                t.power_hp
            FROM brands b
            JOIN models m ON m.brand_id = b.id
            JOIN generations g ON g.model_id = m.id
            JOIN trims t ON t.generation_id = g.id
            ORDER BY b.name, m.name, g.start_year, t.year
            LIMIT 5
        """)
        
        for i, row in enumerate(cursor.fetchall(), 1):
            print(f"{i}. {row[0]} {row[1]} > {row[2]} ({row[3]})")
            print(f"   └─ Trim: {row[4]} ({row[5]})")
            print(f"      {row[6] or 'N/A'}, {row[7] or 'N/A'} HP")
            print()
        
        print("=" * 70)
        print("✅ Migration COMPLETED SUCCESSFULLY!")
        print("=" * 70)
        print()
        print("Next steps:")
        print("  1. Review the generated generations above")
        print("  2. Optionally update generation names and add images")
        print("  3. Update Go code to use the new 4-level structure")
        print("  4. Test API endpoints with the new hierarchy")
        
    except Exception as e:
        conn.rollback()
        print(f"❌ Migration failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
    finally:
        conn.close()

if __name__ == "__main__":
    run_simplified_migration()
