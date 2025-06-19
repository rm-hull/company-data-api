package internal

import (
	"company-data-api/models"
	repo "company-data-api/repositories"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	Results     []models.CompanyDataWithLocation `json:"results"`
	Attribution []string                         `json:"attribution"`
}

type GroupedSearchResponse struct {
	Results     map[string][]models.CompanyDataWithLocation `json:"results"`
	Attribution []string                                    `json:"attribution"`
}

func Search(repo repo.SearchRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		bbox, err := parseBBox(c.Query("bbox"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		results := make([]models.CompanyDataWithLocation, 0, 1000)
		err = repo.Find(bbox, func(companyData *models.CompanyDataWithLocation) {
			results = append(results, *companyData)
		})

		if err != nil {
			log.Printf("error while fetching company data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
			return
		}

		c.JSON(http.StatusOK, SearchResponse{
			Results:     results,
			Attribution: ATTRIBUTION,
		})
	}
}

func GroupByPostcode(repo repo.SearchRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		bbox, err := parseBBox(c.Query("bbox"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		results := make(map[string][]models.CompanyDataWithLocation, 100)
		err = repo.Find(bbox, func(companyData *models.CompanyDataWithLocation) {
			arr, exists := results[companyData.RegAddressPostCode]
			if !exists {
				arr = make([]models.CompanyDataWithLocation, 0, 10)
			}
			results[companyData.RegAddressPostCode] = append(arr, *companyData)
		})

		if err != nil {
			log.Printf("error while fetching company data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
			return
		}

		c.JSON(http.StatusOK, GroupedSearchResponse{
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
