package main

import (
	"errors"
	"fmt"
	"time"
)

type DataSourceType string
type AnalysisType string
type ReportTemplate string

const (
	DataSourceMySQL      DataSourceType = "MySQL"
	DataSourcePostgreSQL DataSourceType = "PostgreSQL"
	DataSourceCSV        DataSourceType = "CSV"
	DataSourceAPI        DataSourceType = "API"

	AnalysisTimeSeries  AnalysisType = "Time Series Analysis"
	AnalysisCategory    AnalysisType = "Category Analysis"
	AnalysisCorrelation AnalysisType = "Correlation Analysis"
	AnalysisPrediction  AnalysisType = "Prediction Analysis"

	TemplateSummary  ReportTemplate = "Summary Report"
	TemplateDetailed ReportTemplate = "Detailed Report"
	TemplateCustom   ReportTemplate = "Custom Report"
)

type DataSource struct {
	ID       string
	Type     DataSourceType
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type AnalysisParameters struct {
	StartDate    time.Time
	EndDate      time.Time
	AnalysisType AnalysisType
	Filters      map[string]string
}

type ReportSection struct {
	Name     string
	Included bool
}

type Report struct {
	ID              string
	Title           string
	DataSource      DataSource
	Parameters      AnalysisParameters
	Template        ReportTemplate
	Sections        []ReportSection
	GeneratedAt     time.Time
	Data            []map[string]string
	Summary         string
	Recommendations []string
}

var (
	ErrDataSourceNil        = errors.New("data source cannot be nil")
	ErrHostRequired         = errors.New("host is required")
	ErrPortRequired         = errors.New("port is required")
	ErrUsernameRequired     = errors.New("username is required")
	ErrPasswordRequired     = errors.New("password is required")
	ErrStartDateRequired    = errors.New("start date is required")
	ErrEndDateRequired      = errors.New("end date is required")
	ErrAnalysisTypeRequired = errors.New("analysis type is required")
	ErrTemplateRequired     = errors.New("report template is required")
)

func getDataSources() []DataSource {
	return []DataSource{
		{
			ID:       "ds_001",
			Type:     DataSourceMySQL,
			Host:     "db.example.com",
			Port:     "3306",
			Username: "analyst",
			Database: "sales_data",
		},
		{
			ID:       "ds_002",
			Type:     DataSourcePostgreSQL,
			Host:     "analytics.example.com",
			Port:     "5432",
			Username: "analyst",
			Database: "customer_data",
		},
		{
			ID:       "ds_003",
			Type:     DataSourceCSV,
			Host:     "files.example.com",
			Port:     "",
			Username: "analyst",
			Database: "exports",
		},
		{
			ID:       "ds_004",
			Type:     DataSourceAPI,
			Host:     "api.example.com",
			Port:     "443",
			Username: "analyst",
			Database: "",
		},
	}
}

func getReportTemplates() []ReportTemplate {
	return []ReportTemplate{
		TemplateSummary,
		TemplateDetailed,
		TemplateCustom,
	}
}

func getReportSections() []ReportSection {
	return []ReportSection{
		{Name: "Overview", Included: true},
		{Name: "Detailed Data", Included: true},
		{Name: "Graphs", Included: true},
		{Name: "Recommendations", Included: true},
	}
}

func generateReport(dataSource DataSource, params AnalysisParameters, template ReportTemplate, sections []ReportSection) (*Report, error) {
	if err := validateDataSource(&dataSource); err != nil {
		return nil, err
	}

	if err := validateParameters(&params); err != nil {
		return nil, err
	}

	if template == "" {
		return nil, ErrTemplateRequired
	}

	// Simulate data retrieval and analysis
	report := &Report{
		ID:              fmt.Sprintf("report_%d", time.Now().UnixNano()),
		Title:           fmt.Sprintf("%s Analysis Report", params.AnalysisType),
		DataSource:      dataSource,
		Parameters:      params,
		Template:        template,
		Sections:        sections,
		GeneratedAt:     time.Now().UTC(),
		Data:            generateSampleData(params),
		Summary:         generateSummary(params),
		Recommendations: generateRecommendations(params),
	}

	return report, nil
}

func validateDataSource(ds *DataSource) error {
	if ds == nil {
		return ErrDataSourceNil
	}
	if ds.Host == "" {
		return ErrHostRequired
	}
	if ds.Port == "" && ds.Type != DataSourceCSV {
		return ErrPortRequired
	}
	if ds.Username == "" {
		return ErrUsernameRequired
	}
	if ds.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func validateParameters(params *AnalysisParameters) error {
	if params.StartDate.IsZero() {
		return ErrStartDateRequired
	}
	if params.EndDate.IsZero() {
		return ErrEndDateRequired
	}
	if params.AnalysisType == "" {
		return ErrAnalysisTypeRequired
	}
	return nil
}

func generateSampleData(params AnalysisParameters) []map[string]string {
	// Generate sample data based on analysis type
	var data []map[string]string

	switch params.AnalysisType {
	case AnalysisTimeSeries:
		// Generate time series data
		currentDate := params.StartDate
		for currentDate.Before(params.EndDate) || currentDate.Equal(params.EndDate) {
			value := 100 + (currentDate.Day() % 20) // Simple variation
			change := "+5%"
			if currentDate.Day()%3 == 0 {
				change = "-2%"
			}

			data = append(data, map[string]string{
				"Date":   currentDate.Format("2006-01-02"),
				"Value":  fmt.Sprintf("%d", value),
				"Change": change,
			})

			currentDate = currentDate.AddDate(0, 0, 1)
		}

	case AnalysisCategory:
		// Generate category data
		categories := []string{"Electronics", "Clothing", "Food", "Books", "Home"}
		for _, category := range categories {
			value := 100 + (len(category) * 10)
			percentage := float64(value) / 500.0 * 100.0

			data = append(data, map[string]string{
				"Category":   category,
				"Value":      fmt.Sprintf("%d", value),
				"Percentage": fmt.Sprintf("%.1f%%", percentage),
			})
		}

	case AnalysisCorrelation:
		// Generate correlation data
		for i := 0; i < 10; i++ {
			x := 10 + i*5
			y := 20 + i*3 + (i%3)*2

			data = append(data, map[string]string{
				"X": fmt.Sprintf("%d", x),
				"Y": fmt.Sprintf("%d", y),
			})
		}

	case AnalysisPrediction:
		// Generate prediction data
		currentDate := params.StartDate
		for i := 0; i < 10; i++ {
			actual := 100 + i*5
			predicted := actual + (i%3)*2 - 1

			data = append(data, map[string]string{
				"Date":       currentDate.Format("2006-01-02"),
				"Actual":     fmt.Sprintf("%d", actual),
				"Predicted":  fmt.Sprintf("%d", predicted),
				"Difference": fmt.Sprintf("%d", predicted-actual),
			})

			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	return data
}

func generateSummary(params AnalysisParameters) string {
	switch params.AnalysisType {
	case AnalysisTimeSeries:
		return "This report analyzes the time series data for the selected period. The data shows a general upward trend with some fluctuations. Key insights include a 15% overall growth and three significant spikes in activity."
	case AnalysisCategory:
		return "This category analysis reveals that Electronics and Clothing are the top performing categories, accounting for 45% of total value. Food and Books show steady performance, while Home category has shown recent decline."
	case AnalysisCorrelation:
		return "The correlation analysis indicates a strong positive relationship between variables X and Y (r=0.85). This suggests that changes in X are likely to cause proportional changes in Y."
	case AnalysisPrediction:
		return "The prediction model shows good accuracy with an average error margin of 3.2%. The model predicts continued growth in the next quarter, with a 12% increase expected."
	default:
		return "This report provides an analysis of the selected data based on the specified parameters."
	}
}

func generateRecommendations(params AnalysisParameters) []string {
	switch params.AnalysisType {
	case AnalysisTimeSeries:
		return []string{
			"Invest in marketing during the identified peak periods",
			"Consider seasonal inventory adjustments based on the pattern",
			"Implement targeted promotions during low activity periods",
		}
	case AnalysisCategory:
		return []string{
			"Increase investment in the Electronics category",
			"Review pricing strategy for the Home category",
			"Explore cross-promotion opportunities between top categories",
		}
	case AnalysisCorrelation:
		return []string{
			"Leverage the strong correlation to optimize resource allocation",
			"Consider A/B testing to validate the relationship",
			"Monitor for any changes in the correlation strength",
		}
	case AnalysisPrediction:
		return []string{
			"Prepare for the predicted growth in the next quarter",
			"Review current capacity to ensure it can handle the expected increase",
			"Consider early procurement to avoid potential supply constraints",
		}
	default:
		return []string{
			"Review the data for any unexpected patterns",
			"Consider additional analysis to gain deeper insights",
			"Share findings with relevant stakeholders",
		}
	}
}
