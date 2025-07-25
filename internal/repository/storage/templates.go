package storage

const (
	GetUserAllInfoTemplate = `SELECT u.id, u.first_name, u.last_name, u.gender, u.age, u.nationality, 
    	ARRAY_AGG(e.email) FILTER (WHERE e.email IS NOT NULL) AS emails 
	FROM Users u LEFT JOIN Emails e ON u.id = e.user_id 
	WHERE u.last_name = $1 GROUP BY u.id;`

	//GetUserAllInfoBySecondNameTemplate = `SELECT u.id, u.first_name, u.last_name, u.gender, u.age, u.nationality,
	//	ARRAY_AGG(e.email) FILTER (WHERE e.email IS NOT NULL) AS emails
	//FROM Users u LEFT JOIN Emails e ON u.id = e.user_id
	//WHERE u.last_name = $1 GROUP BY u.last_name;`

	GetAllUsersTemplate = `SELECT u.id, u.first_name, u.last_name, u.gender, u.age, u.nationality, 
    	ARRAY_AGG(e.email) FILTER (WHERE e.email IS NOT NULL) AS emails 
	FROM Users u LEFT JOIN Emails e ON u.id = e.user_id
	GROUP BY u.id`

	GetAllUserEmailsTemplate = `SELECT id, user_id, email FROM Emails WHERE user_id = $1;`

	GetUserFriendsTemplate = `SELECT 
    	u.id AS friend_id,
    	u.first_name,
    	u.last_name
	FROM Friends f
	JOIN Users u ON 
		(f.id_second_friend = u.id AND f.id_first_friend = $1) OR 
    	(f.id_first_friend = u.id AND f.id_second_friend = $1)
	WHERE $1 IN (f.id_first_friend, f.id_second_friend);`

	AddUserInfoTemplate = `INSERT INTO Users(first_name, last_name, gender, nationality, age) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	AddEmailTemplate = `INSERT INTO Emails(user_id, email) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING;`

	AddFriendshipTemplate = `INSERT INTO Friends(id_first_friend, id_second_friend) VALUES ($1, $2) ON CONFLICT (id_first_friend, id_second_friend) DO NOTHING;`

	UpdateUserInfoTemplate = `UPDATE Users SET first_name = $2, last_name = $3, gender = $4, nationality = $5, age = $6 WHERE id = $1;`

	DeleteUserTemplate = `DELETE FROM Users WHERE id = $1;`

	DeleteEmailTemplate = `DELETE FROM Emails WHERE id = $1;`

	DeleteFriendshipTemplate = `DELETE FROM Friends WHERE id_first_friend = $1 AND id_second_friend = $2;`
)
