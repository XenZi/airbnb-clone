export const formatErrors = (key: string): string => {
  switch (key) {
    case 'username':
      return 'Username must be 3+ characters long, and must be created with letters only. \n';
    case 'password':
      return 'Password must be 8+ characters long, and must contain 1 uppercase, 1 number, 1 special character. \n';
    case 'firstName' || 'lastName':
      return 'First Name and Last Name must be longer than 2 characters. \n';
    default:
      return '';
  }
};
