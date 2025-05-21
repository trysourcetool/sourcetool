export type TicketStatus = string;
export type TicketPriority = string;

export const statusOpen: TicketStatus = 'open' as const;
export const statusInProgress: TicketStatus = 'in_progress' as const;
export const statusResolved: TicketStatus = 'resolved' as const;
export const statusClosed: TicketStatus = 'closed' as const;

export const priorityLow: TicketPriority = 'low' as const;
export const priorityMedium: TicketPriority = 'medium' as const;
export const priorityHigh: TicketPriority = 'high' as const;
export const priorityUrgent: TicketPriority = 'urgent' as const;

export type Ticket = {
  ID: string;
  Title: string;
  Description: string;
  CustomerID: string;
  Status: TicketStatus;
  Priority: TicketPriority;
  Assignee: string;
  CreatedAt: Date;
  UpdatedAt: Date;
};

const generateTestTickets = (): Ticket[] => {
  const baseTime = new Date(2024, 1, 1, 0, 0, 0, 0);
  return [
    {
      ID: '1',
      Title: 'Cannot login to account',
      Description: 'User reported unable to login to their account',
      CustomerID: 'cust_001',
      Status: statusOpen,
      Priority: priorityHigh,
      Assignee: 'support_agent_1',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 0),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 0),
    },
    {
      ID: '2',
      Title: 'Payment processing error',
      Description: 'Payment is not being processed correctly',
      CustomerID: 'cust_002',
      Status: statusInProgress,
      Priority: priorityUrgent,
      Assignee: 'support_agent_2',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 1),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 1),
    },
    {
      ID: '3',
      Title: 'Feature request: Dark mode',
      Description: 'Customer would like to have dark mode option',
      CustomerID: 'cust_003',
      Status: statusOpen,
      Priority: priorityLow,
      Assignee: 'support_agent_1',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 2),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 2),
    },
    {
      ID: '4',
      Title: 'Account deletion request',
      Description:
        'Customer wants to delete their account and all associated data',
      CustomerID: 'cust_004',
      Status: statusInProgress,
      Priority: priorityMedium,
      Assignee: 'support_agent_3',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 3),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 3),
    },
    {
      ID: '5',
      Title: 'Subscription renewal failed',
      Description: 'Automatic renewal failed due to expired credit card',
      CustomerID: 'cust_005',
      Status: statusResolved,
      Priority: priorityHigh,
      Assignee: 'support_agent_2',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 4),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 4),
    },
    {
      ID: '6',
      Title: 'Mobile app crashes on startup',
      Description: 'App crashes immediately after launching on iOS 17.2',
      CustomerID: 'cust_006',
      Status: statusOpen,
      Priority: priorityUrgent,
      Assignee: 'support_agent_4',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 5),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 5),
    },
    {
      ID: '7',
      Title: 'Data export format issue',
      Description: 'CSV export includes incorrect date format',
      CustomerID: 'cust_007',
      Status: statusInProgress,
      Priority: priorityMedium,
      Assignee: 'support_agent_1',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 6),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 6),
    },
    {
      ID: '8',
      Title: 'API rate limiting concerns',
      Description: 'Customer hitting rate limits during peak hours',
      CustomerID: 'cust_008',
      Status: statusOpen,
      Priority: priorityHigh,
      Assignee: 'support_agent_3',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 7),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 7),
    },
    {
      ID: '9',
      Title: 'Billing address update',
      Description: 'Need to update billing address for tax purposes',
      CustomerID: 'cust_009',
      Status: statusResolved,
      Priority: priorityLow,
      Assignee: 'support_agent_2',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 8),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 8),
    },
    {
      ID: '10',
      Title: 'Integration with new CRM',
      Description: 'Request for integration with Salesforce CRM',
      CustomerID: 'cust_010',
      Status: statusOpen,
      Priority: priorityMedium,
      Assignee: 'support_agent_4',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 9),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 9),
    },
    {
      ID: '11',
      Title: 'Password reset not working',
      Description: 'Password reset emails not being received',
      CustomerID: 'cust_011',
      Status: statusInProgress,
      Priority: priorityHigh,
      Assignee: 'support_agent_1',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 10),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 10),
    },
    {
      ID: '12',
      Title: 'Report generation timeout',
      Description: 'Large reports timing out after 5 minutes',
      CustomerID: 'cust_012',
      Status: statusOpen,
      Priority: priorityMedium,
      Assignee: 'support_agent_3',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 11),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 11),
    },
    {
      ID: '13',
      Title: 'Two-factor authentication issues',
      Description: '2FA codes not being accepted',
      CustomerID: 'cust_013',
      Status: statusResolved,
      Priority: priorityHigh,
      Assignee: 'support_agent_2',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 12),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 12),
    },
    {
      ID: '14',
      Title: 'Data import validation errors',
      Description: 'CSV import failing with validation errors',
      CustomerID: 'cust_014',
      Status: statusOpen,
      Priority: priorityMedium,
      Assignee: 'support_agent_4',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 13),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 13),
    },
    {
      ID: '15',
      Title: 'Email notifications delayed',
      Description: 'System emails being delayed by 2-3 hours',
      CustomerID: 'cust_015',
      Status: statusInProgress,
      Priority: priorityHigh,
      Assignee: 'support_agent_1',
      CreatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 14),
      UpdatedAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 14),
    },
  ];
};

const matchesFilter = (
  ticket: Ticket,
  title: string,
  customerId: string,
  status: TicketStatus,
  priority: TicketPriority,
  assignee: string,
) => {
  return (
    (title === '' || ticket.Title === title) &&
    (customerId === '' || ticket.CustomerID === customerId) &&
    (status === '' || ticket.Status === status) &&
    (priority === '' || ticket.Priority === priority) &&
    (assignee === '' || ticket.Assignee === assignee)
  );
};

const filterTickets = (
  tickets: Ticket[],
  title: string,
  customerId: string,
  status: TicketStatus,
  priority: TicketPriority,
  assignee: string,
) => {
  return tickets.filter((ticket) => {
    return matchesFilter(ticket, title, customerId, status, priority, assignee);
  });
};

const validateTicket = (t: Ticket) => {
  if (t.Title === '') {
    return 'Title is required';
  }
  if (t.Description === '') {
    return 'Description is required';
  }
  if (t.CustomerID === '') {
    return 'Customer ID is required';
  }
  if (
    t.Status !== '' &&
    t.Status !== statusOpen &&
    t.Status !== statusInProgress &&
    t.Status !== statusResolved &&
    t.Status !== statusClosed
  ) {
    return 'Invalid status';
  }
  if (
    t.Priority !== '' &&
    t.Priority !== priorityLow &&
    t.Priority !== priorityMedium &&
    t.Priority !== priorityHigh &&
    t.Priority !== priorityUrgent
  ) {
    return 'Invalid priority';
  }
  return null;
};

export const listTickets = (
  title: string,
  customerId: string,
  status: TicketStatus,
  priority: TicketPriority,
  assignee: string,
) => {
  const tickets = generateTestTickets();

  if (
    title === '' &&
    customerId === '' &&
    status === '' &&
    priority === '' &&
    assignee === ''
  ) {
    return tickets;
  }

  return filterTickets(tickets, title, customerId, status, priority, assignee);
};

export const createTicket = (t: Ticket | null) => {
  if (t === null) {
    throw new Error('Ticket cannot be null');
  }

  if (validateTicket(t) !== null) {
    throw new Error('Ticket is invalid');
  }

  if (t.ID === '') {
    t.ID = `ticket_${Date.now()}`;
  }

  const now = new Date();
  if (!t.CreatedAt) {
    t.CreatedAt = now;
  }
  t.UpdatedAt = now;

  return t;
};

export const updateTicket = (t: Ticket | null) => {
  if (t === null) {
    throw new Error('Ticket cannot be null');
  }

  if (validateTicket(t) !== null) {
    throw new Error('Ticket is invalid');
  }

  t.UpdatedAt = new Date();

  return t;
};
