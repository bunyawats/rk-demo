-- name: ListCustomers :many
SELECT fname, lname, age
FROM customer
ORDER BY fname, lname;

-- name: CreateCustomer :execresult
INSERT INTO customer(
    fname, lname, age
) VALUES(
    ?, ?, ?
);

-- name: UpdateCustomer :execresult
UPDATE customer
SET fname=?,  lname=?,  age=?
WHERE cusid=?;

-- name: DeleteCustomer :exec
DELETE FROM customer
WHERE cusid=?;

