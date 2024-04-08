CREATE TABLE cars
(
    car_id SERIAL PRIMARY KEY,
    reg_num VARCHAR(255) NOT NULL UNIQUE,
    mark VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    year INT NOT NULL
);

CREATE INDEX idx_reg_num ON cars(reg_num);

CREATE TABLE owners
(
    owner_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) NOT NULL
);

CREATE TABLE cars_owners
(
    car_id INT NOT NULL,
    owner_id INT NOT NULL,
    PRIMARY KEY (car_id, owner_id),
    FOREIGN KEY (car_id) REFERENCES cars(car_id),
    FOREIGN KEY (owner_id) REFERENCES owners(owner_id)
);
