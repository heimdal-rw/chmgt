CREATE TABLE users (
  id INT PRIMARY KEY,
  username TEXT,
  password TEXT,
  firstname TEXT,
  lastname TEXT,
  email TEXT
);

CREATE TABLE changeRequest (
  id INT PRIMARY KEY,
  title TEXT NOT NULL,
  authorId INT NOT NULL,
  requesterId INT NOT NULL,
  description TEXT NULL,
  reason TEXT NULL,
  risk TEXT NULL,
  steps TEXT NULL,
  revert TEXT NULL,
  FOREIGN KEY (authorId) REFERENCES users(id),
  FOREIGN KEY (requesterId) REFERENCES users(id)
);

CREATE TABLE signoff (
  id INT PRIMARY KEY,
  changeRequestId INT,
  signerId INT,
  signed INT,
  FOREIGN KEY (changeRequestId) REFERENCES changeRequest(id),
  FOREIGN KEY (signerId) REFERENCES users(id)
);

CREATE TABLE configuration (
  id INT PRIMARY KEY,
  key TEXT,
  value TEXT
);

