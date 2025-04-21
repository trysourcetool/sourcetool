type User = {
  id: string;
  name: string;
  email: string;
  age: number;
  gender: string;
  createdAt: Date;
};

const generateUsers = (): User[] => {
  const baseTime = new Date(Date.UTC(2024, 0, 0, 0, 0, 0, 0));
  return [
    {
      id: '1',
      name: 'John Doe 001',
      email: 'john.doe+001@acme.com',
      age: 25,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 0),
    },
    {
      id: '2',
      name: 'John Doe 002',
      email: 'john.doe+002@acme.com',
      age: 30,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 1),
    },
    {
      id: '3',
      name: 'Jane Doe 003',
      email: 'jane.doe+003@acme.com',
      age: 35,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 2),
    },
    {
      id: '4',
      name: 'John Doe 004',
      email: 'john.doe+004@acme.com',
      age: 28,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 3),
    },
    {
      id: '5',
      name: 'Jane Doe 005',
      email: 'jane.doe+005@acme.com',
      age: 32,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 4),
    },
    {
      id: '6',
      name: 'John Doe 006',
      email: 'john.doe+006@acme.com',
      age: 27,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 5),
    },
    {
      id: '7',
      name: 'Jane Doe 007',
      email: 'jane.doe+007@acme.com',
      age: 31,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 6),
    },
    {
      id: '8',
      name: 'John Doe 008',
      email: 'john.doe+008@acme.com',
      age: 29,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 7),
    },
    {
      id: '9',
      name: 'Jane Doe 009',
      email: 'jane.doe+009@acme.com',
      age: 33,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 8),
    },
    {
      id: '10',
      name: 'John Doe 010',
      email: 'john.doe+010@acme.com',
      age: 26,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 9),
    },
    {
      id: '11',
      name: 'Jane Doe 011',
      email: 'jane.doe+011@acme.com',
      age: 34,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 10),
    },
    {
      id: '12',
      name: 'John Doe 012',
      email: 'john.doe+012@acme.com',
      age: 28,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 11),
    },
    {
      id: '13',
      name: 'Jane Doe 013',
      email: 'jane.doe+013@acme.com',
      age: 30,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 12),
    },
    {
      id: '14',
      name: 'John Doe 014',
      email: 'john.doe+014@acme.com',
      age: 32,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 13),
    },
    {
      id: '15',
      name: 'Jane Doe 015',
      email: 'jane.doe+015@acme.com',
      age: 29,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 14),
    },
    {
      id: '16',
      name: 'John Doe 016',
      email: 'john.doe+016@acme.com',
      age: 31,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 15),
    },
    {
      id: '17',
      name: 'Jane Doe 017',
      email: 'jane.doe+017@acme.com',
      age: 27,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 16),
    },
    {
      id: '18',
      name: 'John Doe 018',
      email: 'john.doe+018@acme.com',
      age: 33,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 17),
    },
    {
      id: '19',
      name: 'Jane Doe 019',
      email: 'jane.doe+019@acme.com',
      age: 35,
      gender: 'female',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 18),
    },
    {
      id: '20',
      name: 'John Doe 020',
      email: 'john.doe+020@acme.com',
      age: 30,
      gender: 'male',
      createdAt: new Date(baseTime.getTime() + 24 * 60 * 60 * 1000 * 19),
    },
  ];
};

const matchesFilter = (
  user: User,
  name: string,
  email: string,
  age: number,
  gender: string,
  createdAt?: Date,
): boolean => {
  return (
    (name === '' || user.name === name) &&
    (email === '' || user.email === email) &&
    (age === 0 || user.age === age) &&
    (gender === '' || user.gender === gender) &&
    (!createdAt ||
      createdAt.getTime() === 0 ||
      user.createdAt.getTime() === createdAt.getTime())
  );
};

const filterUsers = (
  users: User[],
  name: string,
  email: string,
  age: number,
  gender: string,
  createdAt?: Date,
): User[] => {
  return users.filter((user) =>
    matchesFilter(user, name, email, age, gender, createdAt),
  );
};

export const listUsers = (
  name: string,
  email: string,
  age: number,
  gender: string,
  createdAt?: Date,
): User[] => {
  const users = generateUsers();

  if (
    name === '' &&
    email === '' &&
    age === 0 &&
    gender === '' &&
    (!createdAt || createdAt.getTime() === 0)
  ) {
    return users;
  }

  return filterUsers(users, name, email, age, gender, createdAt);
};
