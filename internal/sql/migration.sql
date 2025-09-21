CREATE TABLE IF NOT EXISTS code_point (
    post_code TEXT NOT NULL PRIMARY KEY,
    easting NUMERIC NOT NULL,
    northing NUMERIC NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_code_point_easting_northing
ON code_point (easting, northing);


CREATE TABLE IF NOT EXISTS company_data (
    company_name TEXT NOT NULL,
    company_number TEXT NOT NULL PRIMARY KEY,

    reg_address_care_of TEXT,
    reg_address_po_box TEXT,
    reg_address_address_line_1 TEXT NOT NULL,
    reg_address_address_line_2 TEXT,
    reg_address_post_town TEXT NOT NULL,
    reg_address_county TEXT,
    reg_address_country TEXT,
    reg_address_post_code TEXT NOT NULL,

    company_category TEXT NOT NULL,
    company_status TEXT NOT NULL,
    country_of_origin TEXT,

    dissolution_date TIMESTAMP,
    incorporation_date TIMESTAMP NOT NULL,
    accounts_account_ref_day NUMERIC NOT NULL,
    accounts_account_ref_month NUMERIC NOT NULL,

    accounts_next_due_date TIMESTAMP,
    accounts_last_made_up_date TIMESTAMP,
    accounts_account_category TEXT,
    returns_next_due_date TIMESTAMP,
    returns_last_made_up_date TIMESTAMP,

    mortgages_num_charges NUMERIC NOT NULL,
    mortgages_num_outstanding NUMERIC NOT NULL,
    mortgages_num_part_satisfied NUMERIC NOT NULL,
    mortgages_num_satisfied NUMERIC NOT NULL,
    sic_code_1 TEXT NOT NULL,
    sic_code_2 TEXT,
    sic_code_3 TEXT,
    sic_code_4 TEXT,
    limited_partnerships_num_gen_partners NUMERIC NOT NULL,
    limited_partnerships_num_lim_partners NUMERIC NOT NULL,
    uri TEXT NOT NULL,
    conf_stmt_next_due_date TIMESTAMP,
    conf_stmt_last_made_up_date TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_company_data_reg_address_post_code
ON company_data (reg_address_post_code);
