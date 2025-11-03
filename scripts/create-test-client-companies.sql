-- Create Test Client Companies
-- This creates 5 test client companies for development and testing purposes

-- Company 1: TechCorp Solutions
INSERT INTO client_companies (name, contact_info, created_at, updated_at)
VALUES (
    'TechCorp Solutions',
    'Contact: John Smith
Email: contact@techcorp-solutions.com
Phone: +1 (555) 123-4567
Address: 123 Tech Street, San Francisco, CA 94102',
    NOW(),
    NOW()
)
ON CONFLICT ((LOWER(name))) DO UPDATE
SET 
    contact_info = EXCLUDED.contact_info,
    updated_at = NOW();

-- Company 2: Global Innovations Inc
INSERT INTO client_companies (name, contact_info, created_at, updated_at)
VALUES (
    'Global Innovations Inc',
    'Contact: Sarah Johnson
Email: info@globalinnovations.com
Phone: +1 (555) 234-5678
Address: 456 Innovation Blvd, Austin, TX 78701',
    NOW(),
    NOW()
)
ON CONFLICT ((LOWER(name))) DO UPDATE
SET 
    contact_info = EXCLUDED.contact_info,
    updated_at = NOW();

-- Company 3: Digital Dynamics LLC
INSERT INTO client_companies (name, contact_info, created_at, updated_at)
VALUES (
    'Digital Dynamics LLC',
    'Contact: Michael Chen
Email: support@digitaldynamics.com
Phone: +1 (555) 345-6789
Address: 789 Digital Ave, Seattle, WA 98101',
    NOW(),
    NOW()
)
ON CONFLICT ((LOWER(name))) DO UPDATE
SET 
    contact_info = EXCLUDED.contact_info,
    updated_at = NOW();

-- Company 4: CloudFirst Technologies
INSERT INTO client_companies (name, contact_info, created_at, updated_at)
VALUES (
    'CloudFirst Technologies',
    'Contact: Emma Williams
Email: hello@cloudfirst.tech
Phone: +1 (555) 456-7890
Address: 321 Cloud Way, Boston, MA 02101',
    NOW(),
    NOW()
)
ON CONFLICT ((LOWER(name))) DO UPDATE
SET 
    contact_info = EXCLUDED.contact_info,
    updated_at = NOW();

-- Company 5: NextGen Software Group
INSERT INTO client_companies (name, contact_info, created_at, updated_at)
VALUES (
    'NextGen Software Group',
    'Contact: David Rodriguez
Email: contact@nextgensoftware.com
Phone: +1 (555) 567-8901
Address: 555 Future Pkwy, Denver, CO 80202',
    NOW(),
    NOW()
)
ON CONFLICT ((LOWER(name))) DO UPDATE
SET 
    contact_info = EXCLUDED.contact_info,
    updated_at = NOW();

-- Display results
\echo ''
\echo '======================================='
\echo 'Test Client Companies Created!'
\echo '======================================='
\echo ''
\echo 'COMPANY 1: TechCorp Solutions'
\echo '  Contact: John Smith'
\echo '  Email: contact@techcorp-solutions.com'
\echo ''
\echo 'COMPANY 2: Global Innovations Inc'
\echo '  Contact: Sarah Johnson'
\echo '  Email: info@globalinnovations.com'
\echo ''
\echo 'COMPANY 3: Digital Dynamics LLC'
\echo '  Contact: Michael Chen'
\echo '  Email: support@digitaldynamics.com'
\echo ''
\echo 'COMPANY 4: CloudFirst Technologies'
\echo '  Contact: Emma Williams'
\echo '  Email: hello@cloudfirst.tech'
\echo ''
\echo 'COMPANY 5: NextGen Software Group'
\echo '  Contact: David Rodriguez'
\echo '  Email: contact@nextgensoftware.com'
\echo ''
\echo '======================================='
\echo ''

-- Show all test companies
SELECT id, name, 
       SUBSTRING(contact_info FROM 'Email: ([^\n]+)') as email,
       created_at 
FROM client_companies 
WHERE name IN (
    'TechCorp Solutions',
    'Global Innovations Inc',
    'Digital Dynamics LLC',
    'CloudFirst Technologies',
    'NextGen Software Group'
)
ORDER BY name;

