// Where dev defined types most likely should go,
// if they are used in multiple files.

export type Detection = {
  id: string;
  createdBy: string;
  content: string;
  createdAt: string;
  homeworkId: string;
  severity: number;
  filePath: string;
  diffData: string;
};

export type User = {
  email: string;
  token: string;
};

export type Section = {
  id: string;
  name: string;
};

export type Homework = {
  id: string;
  title: string;
  assignedAt: string;
  dueDate: string;
};
