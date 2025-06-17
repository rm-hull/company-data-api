package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	Results     []CompanyDataWithLocation `json:"results"`
	Attribution []string                  `json:"attribution"`
}

type CompanyDataWithLocation struct {
	CompanyData
	Easting  int `json:"easting"`
	Northing int `json:"northing"`
}

const (
	LEFT = iota
	BOTTOM
	RIGHT
	TOP
)

func Search(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		bbox, err := parseBBox(c.Query("bbox"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// In bbox: [LEFT, BOTTOM, RIGHT, TOP]
		// So: bbox[LEFT]=min easting, bbox[BOTTOM]=min northing, bbox[RIGHT]=max easting, bbox[TOP]=max northing
		rows, err := db.Query(`
				SELECT
					cd.company_name, cd.company_number, cd.reg_address_care_of, cd.reg_address_po_box,
					cd.reg_address_address_line_1, cd.reg_address_address_line_2, cd.reg_address_post_town,
					cd.reg_address_county, cd.reg_address_country, cd.reg_address_post_code,
					cd.company_category, cd.company_status, cd.country_of_origin, cd.dissolution_date,
					cd.incorporation_date, cd.accounts_account_ref_day, cd.accounts_account_ref_month,
					cd.accounts_next_due_date, cd.accounts_last_made_up_date, cd.accounts_account_category,
					cd.returns_next_due_date, cd.returns_last_made_up_date, cd.mortgages_num_charges,
					cd.mortgages_num_outstanding, cd.mortgages_num_part_satisfied, cd.mortgages_num_satisfied,
					cd.sic_code_1, cd.sic_code_2, cd.sic_code_3, cd.sic_code_4,
					cd.limited_partnerships_num_gen_partners, cd.limited_partnerships_num_lim_partners,
					cd.uri, cd.conf_stmt_next_due_date, cd.conf_stmt_last_made_up_date,
					cp.easting, cp.northing
				FROM code_point cp
				INNER JOIN company_data cd ON cp.post_code = cd.reg_address_post_code
				WHERE cp.easting BETWEEN ? AND ?
				AND cp.northing BETWEEN ? AND ?
			`,
			bbox[LEFT], bbox[RIGHT], bbox[BOTTOM], bbox[TOP],
		)
		if err != nil {
			log.Printf("error querying database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
			return
		}
		defer func() {
			if err := rows.Close(); err != nil {
				log.Printf("error closing rows: %v", err)
			}
		}()

		var results []CompanyDataWithLocation
		var cd CompanyDataWithLocation

		for rows.Next() {
			if err := rows.Scan(
				&cd.CompanyName,
				&cd.CompanyNumber,
				&cd.RegAddressCareOf,
				&cd.RegAddressPOBox,
				&cd.RegAddressAddressLine1,
				&cd.RegAddressAddressLine2,
				&cd.RegAddressPostTown,
				&cd.RegAddressCounty,
				&cd.RegAddressCountry,
				&cd.RegAddressPostCode,
				&cd.CompanyCategory,
				&cd.CompanyStatus,
				&cd.CountryOfOrigin,
				&cd.DissolutionDate,
				&cd.IncorporationDate,
				&cd.AccountsAccountRefDay,
				&cd.AccountsAccountRefMonth,
				&cd.AccountsNextDueDate,
				&cd.AccountsLastMadeUpDate,
				&cd.AccountsAccountCategory,
				&cd.ReturnsNextDueDate,
				&cd.ReturnsLastMadeUpDate,
				&cd.MortgagesNumMortCharges,
				&cd.MortgagesNumMortOutstanding,
				&cd.MortgagesNumMortPartSatisfied,
				&cd.MortgagesNumMortSatisfied,
				&cd.SICCodeSicText_1,
				&cd.SICCodeSicText_2,
				&cd.SICCodeSicText_3,
				&cd.SICCodeSicText_4,
				&cd.LimitedPartnershipsNumGenPartners,
				&cd.LimitedPartnershipsNumLimPartners,
				&cd.URI,
				&cd.ConfStmtNextDueDate,
				&cd.ConfStmtLastMadeUpDate,
				&cd.Easting,
				&cd.Northing,
			); err != nil {
				log.Printf("error scanning row: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
				return
			}

			results = append(results, cd)
		}
		if err = rows.Err(); err != nil {
			log.Printf("error during rows iteration: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
			return
		}

		c.JSON(http.StatusOK, SearchResponse{
			Results:     results,
			Attribution: ATTRIBUTION,
		})
	}
}

func parseBBox(bboxStr string) ([]float64, error) {
	bboxParts := strings.Split(bboxStr, ",")
	if len(bboxParts) != 4 {
		return nil, fmt.Errorf("bbox must have 4 comma-separated values")
	}

	bbox := make([]float64, 4)
	for i, part := range bboxParts {
		val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid bbox value '%s': not a valid float", part)
		}
		bbox[i] = val
	}

	return bbox, nil
}
