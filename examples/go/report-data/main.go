package main

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/multiselect"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func createReportPage(ui sourcetool.UIBuilder) error {
	// Data Source Selection
	form, submitted := ui.Form("Generate Report", form.WithClearOnSubmit(true))
	form.Markdown("### 1. Select data source")
	dataSources := getDataSources()
	var dataSourceOptions []string
	for _, ds := range dataSources {
		dataSourceOptions = append(dataSourceOptions, string(ds.Type))
	}

	cols := form.Columns(2)
	selectedDataSourceType := cols[0].MultiSelect(
		"Data Source Type",
		multiselect.WithOptions(dataSourceOptions...),
	)

	var selectedDataSource DataSource
	for _, ds := range dataSources {
		if selectedDataSourceType != nil {
			if slices.Contains(selectedDataSourceType.Values, string(ds.Type)) {
				selectedDataSource = ds
				break
			}
		}
	}

	// Data Source Configuration
	form.Markdown("### 2. Configure data source")
	configCols := form.Columns(2)
	host := configCols[0].TextInput("Host", textinput.WithDefaultValue(selectedDataSource.Host))
	port := configCols[1].TextInput("Port", textinput.WithDefaultValue(selectedDataSource.Port))

	authCols := form.Columns(2)
	username := authCols[0].TextInput("Username", textinput.WithDefaultValue(selectedDataSource.Username))
	password := authCols[1].TextInput("Password", textinput.WithPlaceholder("Enter password"))

	// Analysis Parameters
	form.Markdown("### 3. Analysis parameters")
	paramCols := form.Columns(2)
	startDate := paramCols[0].DateInput("Start Date")
	endDate := paramCols[1].DateInput("End Date")

	analysisTypes := []string{
		string(AnalysisTimeSeries),
		string(AnalysisCategory),
		string(AnalysisCorrelation),
		string(AnalysisPrediction),
	}
	analysisType := form.Selectbox(
		"Analysis Type",
		selectbox.WithOptions(analysisTypes...),
	)

	// Report Configuration
	form.Markdown("### 4. Report configuration")
	templates := getReportTemplates()
	var templateOptions []string
	for _, t := range templates {
		templateOptions = append(templateOptions, string(t))
	}
	selectedTemplate := form.Selectbox(
		"Report Template",
		selectbox.WithOptions(templateOptions...),
	)

	// Report Sections
	form.Markdown("### 5. Report sections")
	sections := getReportSections()
	var sectionOptions []string
	for _, s := range sections {
		sectionOptions = append(sectionOptions, s.Name)
	}
	selectedSections := form.MultiSelect(
		"Included Sections",
		multiselect.WithOptions(sectionOptions...),
	)

	if submitted {
		// Update data source with form values
		selectedDataSource.Host = host
		selectedDataSource.Port = port
		selectedDataSource.Username = username
		selectedDataSource.Password = password

		// Create analysis parameters
		params := AnalysisParameters{
			StartDate:    *startDate,
			EndDate:      *endDate,
			AnalysisType: AnalysisType(analysisType.Value),
			Filters:      make(map[string]string),
		}

		// Create report sections
		var reportSections []ReportSection
		for _, sectionName := range selectedSections.Values {
			reportSections = append(reportSections, ReportSection{
				Name:     sectionName,
				Included: true,
			})
		}

		// Generate report
		report, err := generateReport(
			selectedDataSource,
			params,
			ReportTemplate(selectedTemplate.Value),
			reportSections,
		)
		if err != nil {
			return fmt.Errorf("failed to generate report: %v", err)
		}

		// Display report
		ui.Markdown(fmt.Sprintf("## Report: %s", report.Title))
		ui.Markdown(fmt.Sprintf("Generated at: %s", report.GeneratedAt.Format(time.RFC1123)))

		ui.Markdown("### Summary")
		ui.Markdown(report.Summary)

		ui.Markdown("### Data")
		ui.Table(
			report.Data,
			table.WithHeight(10),
		)

		ui.Markdown("### Recommendations")
		for _, rec := range report.Recommendations {
			ui.Markdown(fmt.Sprintf("- %s", rec))
		}
	}

	return nil
}

func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "your_api_key",
		Endpoint: "ws://localhost:3000",
	})

	s.Page("/reports/new", "Create Report", createReportPage)

	if err := s.Listen(); err != nil {
		log.Printf("Failed to listen sourcetool: %v", err)
		s.Close()
		return
	}
}
