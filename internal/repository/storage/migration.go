package storage

const (
	createUsersTableTemplate = `CREATE TABLE IF NOT EXISTS Users(
		id serial primary key,
		first_name text not null,
		last_name text not null,
		gender text not null,
		nationality text not null,
		age integer not null
	);`

	createEmailsTableTemplate = `CREATE TABLE IF NOT EXISTS Emails(
    	id serial primary key,
    	user_id integer not null,
    	email text not null UNIQUE,
                                
        FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE ON UPDATE CASCADE
	);`

	createFriendsTableTemplate = `CREATE TABLE IF NOT EXISTS Friends(
    	id_first_friend integer not null,
    	id_second_friend integer not null,
    	
    	PRIMARY KEY (id_first_friend, id_second_friend),
    
    	FOREIGN KEY (id_first_friend) REFERENCES Users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    
    	FOREIGN KEY (id_second_friend) REFERENCES Users(id) ON DELETE CASCADE ON UPDATE CASCADE
	);`

	createFriendsIndexTemplates = `CREATE INDEX IF NOT EXISTS id_second_first_friend ON Friends(id_second_friend, id_first_friend);`
)
