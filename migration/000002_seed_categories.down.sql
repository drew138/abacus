DELETE FROM categories WHERE created_at = '2026-11-04T00:00:00Z' AND name NOT IN (SELECT category FROM movements) AND name NOT IN (SELECT category FROM merchant_rules);
