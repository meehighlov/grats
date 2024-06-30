package db

import (
	"log"
)

// todo add metada for error logging

func (user *User) Save() error {
	stmt, err := Client.Prepare(
		"INSERT INTO user(id, tgid, name, tgusername, chatid, birthday, isadmin, createdat, updatedat) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) " +
			"ON CONFLICT(tgid) DO UPDATE SET name=$3, tgusername=$4, chatid=$5, birthday=$6, isadmin=$7, updatedat=$9 " +
			"RETURNING id;",
	)
	if err != nil {
		log.Println("Error when trying to prepare statement for saving user: " + err.Error())
		return err
	}
	defer stmt.Close()

    _, _, _ = user.RefresTimestamps()

	insertErr := stmt.QueryRow(
        &user.ID,
        &user.TGId,
        &user.Name,
        &user.TGusername,
        &user.ChatId,
        &user.Birthday,
        &user.IsAdmin,
        &user.CreatedAt,
        &user.UpdatedAt,
    ).Scan(&user.ID)
	if insertErr != nil {
		log.Println("Error when trying to save user: ", insertErr.Error())
		return insertErr
	}
	log.Println("User created/updated")

	return nil
}

func (user *User) Get() error {
	stmt, err := Client.Prepare("SELECT id, tgid, name, tgusername, chatid, birthday, isadmin, createdat, updatedat FROM user WHERE id=$1;")
	if err != nil {
		log.Println("Error when trying to prepare statement for getting user by id: " + err.Error())
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRow(&user.ID)
	if err := result.Scan(
        &user.ID,
        &user.TGId,
        &user.Name,
        &user.TGusername,
        &user.ChatId,
        &user.Birthday,
        &user.IsAdmin,
        &user.CreatedAt,
        &user.UpdatedAt,
    ); err != nil {
		log.Println("Error when trying to get User by ID: " + err.Error())
		return err
	}

	return nil
}

func (user *User) IsExist() bool {
	stmt, err := Client.Prepare("SELECT COUNT(user.id) FROM user WHERE id=$1;")
	if err != nil {
		log.Println("Error when trying to prepare statement for getting user by id: " + err.Error())
		return false
	}
	defer stmt.Close()

	result := stmt.QueryRow(&user.ID)
	var count *int
	if err := result.Scan(&count); err != nil {
		log.Println("Error when trying to get User by ID: " + err.Error())
		return false
	}

	return *count == 1
}

func (friend *Friend) FilterByNotifyDate(date string) ([]Friend, error) {
	stmt, err := Client.Prepare(
		"SELECT id, name, birthday, userid, chatid, notifyat, createdat, updatedat FROM friend WHERE notifyat=$1;",
	)
	if err != nil {
		log.Println("Error when trying to prepare statement for fetching friends: " + err.Error())
		return nil, err
	}
    defer stmt.Close()

	results, err := stmt.Query(date)
	if err != nil {
		log.Println("Error when fetching friends by birthday", err.Error())
		return nil, err
	}
    defer results.Close()

    friends := []Friend{}

	for results.Next() {
		friend := Friend{}
		err := results.Scan(
            &friend.ID,
            &friend.Name,
            &friend.BirthDay,
            &friend.UserId,
            &friend.ChatId,
            friend.GetNotifyAt(),
            &friend.CreatedAt,
            &friend.UpdatedAt,
        )
		if err != nil {
			log.Println("Error when fetching friends by birthday:", err)
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (friend *Friend) Save() error {
	stmt, err := Client.Prepare(
        "INSERT INTO friend(id, name, birthday, userid, chatid, notifyat, createdat, updatedat) " +
        "VALUES($1, $2, $3, $4, $5, $6, $7, $8) " +
        "ON CONFLICT(id) DO UPDATE SET name=$2, birthday=$3, userid=$4, chatid=$5, notifyat=$6, createdat=$7, updatedat=$8 " +
        "RETURNING id;",
    )
	if err != nil {
		log.Println("Error when trying to prepare statement: " + err.Error())
		return err
	}
	defer stmt.Close()

    _, _, _ = friend.RefresTimestamps()

	insertErr := stmt.QueryRow(
        friend.ID,
        friend.Name,
        friend.BirthDay,
        friend.UserId,
        friend.ChatId,
        *friend.GetNotifyAt(),
        friend.CreatedAt,
        friend.UpdatedAt,
    ).Scan(&friend.ID)
	if insertErr != nil {
		log.Printf("Error when trying to save friend: " + insertErr.Error())
		return insertErr
	}
	log.Println("Friend added/updated")

	return nil
}

func (friend *Friend) FilterByUserId() ([]Friend, error) {
	stmt, err := Client.Prepare("SELECT id, name, birthday, userid, chatid, notifyat, createdat, updatedat FROM friend WHERE userid=$1;")
	if err != nil {
		log.Printf("Error when trying to prepare statement for fetching friends for user %s", err.Error())
		return nil, err
	}
    defer stmt.Close()

	results, err := stmt.Query(&friend.UserId)
	if err != nil {
		log.Printf("Error when fetching friends for user with id %d, error: %s", friend.UserId, err.Error())
		return nil, err
	}
    defer results.Close()

	friends := []Friend{}

	for results.Next() {
		friend := Friend{}
		err := results.Scan(
            &friend.ID,
            &friend.Name,
            &friend.BirthDay,
            &friend.UserId,
            &friend.ChatId,
            friend.GetNotifyAt(),
            &friend.CreatedAt,
            &friend.UpdatedAt,
        )
		if err != nil {
			log.Printf("Error when fetching friends for user with id %d, error: %s", friend.UserId, err.Error())
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (access *Access) All() (*map[string]Access, error) {
    stmt, err := Client.Prepare("SELECT id, tgusername FROM access;")
	if err != nil {
		log.Printf("Error when trying to prepare statement for fetching access list, error: %s", err.Error())
		return nil, err
	}
    defer stmt.Close()

	results, err := stmt.Query()
	if err != nil {
		log.Printf("Error when fetching access list, error: %s", err.Error())
		return nil, err
	}
    defer results.Close()

	accessList := make(map[string]Access)

	for results.Next() {
		access := Access{}
		err := results.Scan(&access.ID, &access.TGusername)
		if err != nil {
			log.Printf("Error when fetching access, error: %s", err.Error())
			continue
		}
		accessList[access.TGusername] = access
	}

	return &accessList, nil
}

func (access *Access) Save() error {
	stmt, err := Client.Prepare(
		"INSERT INTO access(id, tgusername, createdat, updatedat) " +
			"VALUES($1, $2, $3, $4) " +
			"ON CONFLICT(tgusername) DO UPDATE SET tgusername=$2, updatedat=$4 " +
			"RETURNING id;",
	)
	if err != nil {
		log.Println("Error when trying to prepare statement for saving access: " + err.Error())
		return err
	}
	defer stmt.Close()

    _, _, _ = access.RefresTimestamps()

	insertErr := stmt.QueryRow(&access.ID, &access.TGusername, &access.CreatedAt, &access.UpdatedAt).Scan(&access.ID)
	if insertErr != nil {
		log.Println("Error when trying to save access: ", insertErr.Error())
		return insertErr
	}
	log.Println("Access created/updated")

	return nil
}

func (access *Access) IsExist() bool {
	stmt, err := Client.Prepare("SELECT COUNT(id) FROM access WHERE tgusername=$1;")
	if err != nil {
		log.Println("Error when trying to prepare statement for getting acesss: " + err.Error())
		return false
	}
	defer stmt.Close()

	result := stmt.QueryRow(&access.TGusername)
	var count *int
	if err := result.Scan(&count); err != nil {
		log.Println("Error when trying to get access: " + err.Error())
		return false
	}

	return *count == 1
}

func (access *Access) Delete() error {
    stmt, err := Client.Prepare(
		`DELETE FROM access WHERE tgusername = $1;`,
	)
	if err != nil {
		log.Println("Error when trying to prepare statement for deleting access: " + err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&access.TGusername)
	if err != nil {
		log.Println("Error when trying to delete access row: " + err.Error())
		return err
	}

	log.Println("Access row deleted")

	return nil
}
