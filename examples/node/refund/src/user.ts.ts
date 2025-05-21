type User = {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
};

export type RefundRequest = {
  userId: string;
  amount: number;
  reason: string;
};

export const listUsers = (name: string, email: string): User[] => {
  const baseTime = new Date(2024, 1, 1, 0, 0, 0, 0);
  const users: User[] = [
    {
      id: '1',
      name: 'John Doe',
      email: 'john.doe@acme.com',
      createdAt: baseTime,
    },
    {
      id: '2',
      name: 'Jane Smith',
      email: 'jane.smith@acme.com',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000),
    },
    {
      id: '3',
      name: 'Bob Lee',
      email: 'bob.lee@acme.com',
      createdAt: new Date(baseTime.getTime() + 48 * 60 * 60 * 1000),
    },
  ];
  const filtered: User[] = [];
  for (const user of users) {
    if (
      (name === '' || user.name === name) &&
      (email === '' || user.email === email)
    ) {
      filtered.push(user);
    }
  }
  return filtered;
};

export const refundStripe = (req: RefundRequest): void => {
  console.log(
    `Refund: userID=${req.userId}, amount=${req.amount}, reason=${req.reason}`,
  );
};
