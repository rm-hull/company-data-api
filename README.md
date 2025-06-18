# Companies Data API

## Preparation

Download .zip files as follows

-   company data from: https://download.companieshouse.gov.uk/en_output.html
-   code point from: https://osdatahub.os.uk/downloads/open/CodePointOpen

Put both .zip files in ./data folder - there is no need to decompress the files - this will be done
automatically as part of the import process. You will need approximately 3.5Gb of free storage to load
all the company data and code points.

## Misc

```sql
explain query plan
select cp.*, cd.* from code_point cp
inner join company_data cd on cp.post_code = cd.reg_address_post_code
where cp.easting between 425000 and 435000
and cp.northing between 450000 and 460000;
```
