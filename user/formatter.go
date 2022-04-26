package user

type UserFormatter struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{}
	formatter.Id = user.ID
	formatter.Name = user.Name
	formatter.Email = user.Email
	formatter.Occupation = user.Occupation
	formatter.Token = token
	return formatter
}
