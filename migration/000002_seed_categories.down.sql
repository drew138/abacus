DELETE FROM categories
WHERE name IN ('Savings', 'Restaurants and cafés', 'Food delivery', 'Groceries', 'Snacks and candy', 'Transfers to people',
               'Rides and transport', 'Pharmacy and health', 'Gym and fitness', 'Online shopping', 'Travel',
               'Phone, internet and subscriptions', 'Gaming', 'Entertainment', 'Other')
  AND name NOT IN (SELECT category FROM movements)
  AND name NOT IN (SELECT category FROM merchant_rules);
