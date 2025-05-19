import { Sourcetool, SourcetoolConfig, UIBuilderType } from '@sourcetool/node';
import {
  analysisCategory,
  analysisCorrelation,
  AnalysisParameters,
  analysisPrediction,
  analysisTimeSeries,
  DataSource,
  generateReport,
  getDataSources,
  getReportSections,
  getReportTemplates,
  ReportSection,
} from './report';

const createReportPage = async (ui: UIBuilderType) => {
  // Data Source Selection
  const [form, submitted] = ui.form('Generate Report', { clearOnSubmit: true });
  form.markdown('### 1. Select data source');
  const dataSources = getDataSources();
  const dataSourceOptions: string[] = [];
  for (const ds of dataSources) {
    dataSourceOptions.push(String(ds.type));
  }

  const cols = form.columns(2);
  const selectedDataSourceType = cols[0].multiSelect('Data Source Type', {
    options: dataSourceOptions,
  });

  const selectedDataSource = dataSources.find((ds) =>
    selectedDataSourceType?.values.includes(String(ds.type)),
  );

  // Data Source Configuration
  form.markdown('### 2. Configure data source');
  const configCols = form.columns(2);
  const host = configCols[0].textInput('Host', {
    defaultValue: selectedDataSource?.host,
  });
  const port = configCols[1].textInput('Port', {
    defaultValue: selectedDataSource?.port,
  });

  const authCols = form.columns(2);
  const username = authCols[0].textInput('Username', {
    defaultValue: selectedDataSource?.username,
  });
  const password = authCols[1].textInput('Password', {
    placeholder: 'Enter password',
  });

  // Analysis Parameters
  form.markdown('### 3. Analysis parameters');
  const paramCols = form.columns(2);
  const startDate = paramCols[0].dateInput('Start Date');
  const endDate = paramCols[1].dateInput('End Date');

  const analysisTypes: string[] = [
    analysisTimeSeries,
    analysisCategory,
    analysisCorrelation,
    analysisPrediction,
  ];
  const analysisType = form.selectbox('Analysis Type', {
    options: analysisTypes,
  });

  // Report Configuration
  form.markdown('### 4. Report configuration');
  const templateOptions = getReportTemplates();

  const selectedTemplate = form.selectbox('Report Template', {
    options: templateOptions,
  });

  // Report Sections
  form.markdown('### 5. Report sections');
  const sectionOptions: string[] = getReportSections().map((s) => s.name);

  const selectedSections = form.multiSelect('Included Sections', {
    options: sectionOptions,
  });

  if (submitted) {
    // Update data source with form values
    if (!selectedDataSource) return;
    selectedDataSource.host = host;
    selectedDataSource.port = port;
    selectedDataSource.username = username;
    selectedDataSource.password = password;

    // Create analysis parameters
    const params: AnalysisParameters = {
      startDate: startDate,
      endDate: endDate,
      analysisType: analysisType?.value ?? '',
      filters: {},
    };

    // Create report sections
    const reportSections: ReportSection[] =
      selectedSections?.values.map((name) => ({
        name,
        included: true,
      })) ?? [];

    // Generate report
    const report = generateReport(
      selectedDataSource as DataSource,
      params,
      selectedTemplate?.value ?? '',
      reportSections,
    );
    if (!report) return;

    // Display report
    ui.markdown(`## Report: ${report.title}`);
    ui.markdown(`Generated at: ${report.generatedAt.toISOString()}`);

    ui.markdown('### Summary');
    ui.markdown(report.summary);

    ui.markdown('### Data');
    ui.table(report.data, { height: 10 });

    ui.markdown('### Recommendations');
    for (const rec of report.recommendations) {
      ui.markdown(`- ${rec}`);
    }
  }
};

const config: SourcetoolConfig = {
  apiKey: 'your_api_key',
  endpoint: 'ws://localhost:3000',
};

const sourcetool = new Sourcetool(config);

sourcetool.page('/reports/new', 'Create Report', createReportPage);

sourcetool.listen();
