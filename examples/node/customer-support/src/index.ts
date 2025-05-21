import {
  TableOnSelect,
  Sourcetool,
  SourcetoolConfig,
  UIBuilder,
} from '@sourcetool/node';
import {
  createTicket,
  listTickets,
  priorityHigh,
  priorityLow,
  priorityMedium,
  priorityUrgent,
  statusClosed,
  statusInProgress,
  statusOpen,
  statusResolved,
  Ticket,
  TicketPriority,
  TicketStatus,
  updateTicket,
} from './ticket';

const listTicketsPage = async (ui: UIBuilder) => {
  const searchCols = ui.columns(3);
  const title = searchCols[0].textInput('Title', {
    placeholder: 'Enter title to filter',
  });
  const customerID = searchCols[1].textInput('Customer ID', {
    placeholder: 'Enter customer ID to filter',
  });
  const assignee = searchCols[2].textInput('Assignee', {
    placeholder: 'Enter assignee to filter',
  });

  const statusCols = ui.columns(2);
  const status = statusCols[0].selectbox('Status', {
    options: [statusOpen, statusInProgress, statusResolved, statusClosed],
  });
  const priority = statusCols[1].selectbox('Priority', {
    options: [priorityLow, priorityMedium, priorityHigh, priorityUrgent],
  });

  const statusValue = status?.value ?? '';
  const priorityValue = priority?.value ?? '';

  const tickets = listTickets(
    title,
    customerID,
    statusValue,
    priorityValue,
    assignee,
  );

  const baseCols = ui.columns(2, { weight: [3, 1] });
  const ticketTable = baseCols[0].table(tickets, {
    height: 10,
    columnOrder: [
      'ID',
      'Title',
      'CustomerID',
      'Status',
      'Priority',
      'Assignee',
      'CreatedAt',
    ],
    onSelect: TableOnSelect.Rerun,
  });

  let defaultTitle: string = '';
  let defaultDescription: string = '';
  let defaultCustomerID: string = '';
  let defaultAssignee: string = '';
  let defaultStatus: TicketStatus = '';
  let defaultPriority: TicketPriority = '';

  if (ticketTable?.selection && ticketTable.selection.row < tickets.length) {
    const selectedData = tickets[ticketTable.selection.row];
    defaultTitle = selectedData.Title;
    defaultDescription = selectedData.Description;
    defaultCustomerID = selectedData.CustomerID;
    defaultStatus = selectedData.Status;
    defaultPriority = selectedData.Priority;
    defaultAssignee = selectedData.Assignee;
  }

  const [form, submitted] = baseCols[1].form('Update Ticket', {
    clearOnSubmit: true,
  });
  const formTitle = form.textInput('Title', {
    placeholder: 'Enter ticket title',
    defaultValue: defaultTitle,
    required: true,
  });
  const formDescription = form.textInput('Description', {
    placeholder: 'Enter ticket description',
    defaultValue: defaultDescription,
    required: true,
  });
  const formCustomerID = form.textInput('Customer ID', {
    placeholder: 'Enter customer ID',
    defaultValue: defaultCustomerID,
    required: true,
  });
  const formStatus = form.selectbox('Status', {
    options: [statusOpen, statusInProgress, statusResolved, statusClosed],
    defaultValue: defaultStatus,
  });
  const formPriority = form.selectbox('Priority', {
    options: [priorityLow, priorityMedium, priorityHigh, priorityUrgent],
    defaultValue: defaultPriority,
  });
  const formAssignee = form.textInput('Assignee', {
    placeholder: 'Enter assignee',
    defaultValue: defaultAssignee,
  });

  if (submitted) {
    // Use the form status and priority values directly
    const formStatusValue = formStatus?.value ?? '';
    const formPriorityValue = formPriority?.value ?? '';

    const ticket: Ticket = {
      ID: '30',
      Title: formTitle,
      Description: formDescription,
      CustomerID: formCustomerID,
      Status: formStatusValue,
      Priority: formPriorityValue,
      Assignee: formAssignee,
      CreatedAt: new Date(),
      UpdatedAt: new Date(),
    };
    const updatedTicket = updateTicket(ticket);
    if (updatedTicket === null) {
      throw new Error('Ticket is invalid');
    }
    ui.markdown('Ticket updated successfully');
  }
};

const createTicketPage = async (ui: UIBuilder) => {
  const [form, submitted] = ui.form('Create Ticket', {
    clearOnSubmit: true,
  });
  const formTitle = form.textInput('Title', {
    placeholder: 'Enter ticket title',
    required: true,
  });
  const formDescription = form.textInput('Description', {
    placeholder: 'Enter ticket description',
    required: true,
  });
  const formCustomerID = form.textInput('Customer ID', {
    placeholder: 'Enter customer ID',
    required: true,
  });
  const formStatus = form.selectbox('Status', {
    options: [statusOpen, statusInProgress, statusResolved, statusClosed],
    defaultValue: statusOpen,
  });
  const formPriority = form.selectbox('Priority', {
    options: [priorityLow, priorityMedium, priorityHigh, priorityUrgent],
    defaultValue: priorityMedium,
  });
  const formAssignee = form.textInput('Assignee', {
    placeholder: 'Enter assignee',
  });

  if (submitted) {
    // Use the form status and priority values directly
    const formStatusValue = formStatus?.value ?? '';
    const formPriorityValue = formPriority?.value ?? '';

    const ticket: Ticket = {
      ID: '30',
      Title: formTitle,
      Description: formDescription,
      CustomerID: formCustomerID,
      Status: formStatusValue,
      Priority: formPriorityValue,
      Assignee: formAssignee,
      CreatedAt: new Date(),
      UpdatedAt: new Date(),
    };
    const createdTicket = createTicket(ticket);
    if (createdTicket === null) {
      throw new Error('Ticket is invalid');
    }
    ui.markdown(`Ticket created successfully with ID: ${ticket.ID}`);
  }
};

const config: SourcetoolConfig = {
  apiKey: 'your_api_key',
  endpoint: 'ws://localhost:3000',
};

const sourcetool = new Sourcetool(config);

sourcetool.page('/tickets', 'Tickets', listTicketsPage);
sourcetool.page('/tickets/new', 'Create ticket', createTicketPage);

sourcetool.listen();
