package main

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoanParams struct {
	Plafon    float64 `json:"plafon"`
	Term      int     `json:"term"`
	Rate      float64 `json:"rate"`
	StartDate string  `json:"start_date"`
}

type Installment struct {
	Number               int     `json:"number"`
	Date                 string  `json:"date"`
	TotalInstallment     float64 `json:"total_installment"`
	PrincipalInstallment float64 `json:"principal_installment"`
	InterestInstallment  float64 `json:"interest_installment"`
	RemainingPrincipal   float64 `json:"remaining_principal"`
}

func roundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func calculateInstallments(params LoanParams) ([]Installment, error) {
	rateMonthly := params.Rate / 12 / 100
	startDate, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		return nil, err
	}
	annuity := params.Plafon * rateMonthly * math.Pow(1+rateMonthly, float64(params.Term)) / (math.Pow(1+rateMonthly, float64(params.Term)) - 1)

	installments := make([]Installment, params.Term)
	remainingPrincipal := params.Plafon

	for i := 0; i < params.Term; i++ {
		interestInstallment := remainingPrincipal * rateMonthly
		principalInstallment := annuity - interestInstallment
		remainingPrincipal -= principalInstallment

		installments[i] = Installment{
			Number:               i + 1,
			Date:                 startDate.AddDate(0, i, 0).Format("2006-01-02"),
			TotalInstallment:     roundFloat(annuity, 2),
			PrincipalInstallment: roundFloat(principalInstallment, 2),
			InterestInstallment:  roundFloat(interestInstallment, 2),
			RemainingPrincipal:   roundFloat(remainingPrincipal, 2),
		}
	}
	return installments, nil
}

func main() {
	r := gin.Default()

	r.POST("/calculate", func(c *gin.Context) {
		var params LoanParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		installments, err := calculateInstallments(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, installments)
	})

	r.Run()
}
