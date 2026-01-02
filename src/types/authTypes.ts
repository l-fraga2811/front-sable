export type LoginFormData = {
  email: string;
  password: string;
  phone?: string | undefined;
};

export type RegisterFormData = {
  username: string;
  email: string;
  password: string;
  phone?: string | undefined;
  confirmPassword: string;
};
