CREATE TABLE users (
  id INT NOT NULL,
  firstname TEXT,
  lastname TEXT,
  email TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE changeRequest (
  id INT NOT NULL,
  title VARCHAR(100) NOT NULL,
  authorId INT NOT NULL,
  requesterId INT NOT NULL,
  description VARCHAR(255) NULL,
  reason VARCHAR(255) NULL,
  risk VARCHAR(45) NULL,
  steps LONGTEXT NULL,
  revert LONGTEXT NULL,
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
)