type User = {
  id: string;
  name: string;
  email: string;
  age: number;
  gender: string;
  createdAt: Date;
};

const validateUser = (user: User): null => {
  if (user.name === '') {
    throw new Error('Name is required');
  }
  if (user.email === '') {
    throw new Error('Email is required');
  }
  if (user.age <= 0) {
    throw new Error('Age must be positive');
  }
  if (user.gender === '') {
    throw new Error('Gender is required');
  }
  if (user.gender !== 'male' && user.gender !== 'female') {
    throw new Error('Gender must be either male or female');
  }
  return null;
};

export const createUser = (user: User): User | null => {
  const u = { ...user };

  const valid = validateUser(u);
  if (valid || u.id === '') {
    return null;
  }

  if (!u.createdAt || u.createdAt.getTime() === 0) {
    u.createdAt = new Date();
  }

  return u;
};
