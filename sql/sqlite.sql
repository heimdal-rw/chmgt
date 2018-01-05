CREATE TABLE users (
  id INT NOT NULL,
  username TEXT,
  password TEXT,
  firstname TEXT,
  lastname TEXT,
  email TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE changeRequest (
  id INT NOT NULL,
  title TEXT NOT NULL,
  authorId INT NOT NULL,
  requesterId INT NOT NULL,
  description TEXT NULL,
  reason TEXT NULL,
  risk TEXT NULL,
  steps TEXT NULL,
  revert TEXT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (authorId) REFERENCES users(id),
  FOREIGN KEY (requesterId) REFERENCES users(id)
);

CREATE TABLE signoff (
  id INT NOT NULL,
  changeRequestId INT,
  signerId INT,
  signed INT,
  PRIMARY KEY (id),
  FOREIGN KEY (changeRequestId) REFERENCES changeRequest(id),
  FOREIGN KEY (signerId) REFERENCES users(id)
);

CREATE TABLE configuration (
  id INT NOT NULL,
  key TEXT,
  value TEXT,
  PRIMARY KEY (id)
);

