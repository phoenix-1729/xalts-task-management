CREATE TABLE approvals (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL,
    approver_id INT NOT NULL,
    comment TEXT,
    approved BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (approver_id) REFERENCES users(id)
);