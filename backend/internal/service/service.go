package service

type DBStorage interface {
	CreateUser(email, password, nickname string) error
	LoginUser(email, password string) (string, error)
	GetUserFromToken(token string) (User, error)

	StartSession(userID int64, mode string) (Session, Question, error)
	SubmitAnswer(userID, sessionID int64, req SubmitAnswer) (AnswerResult, error)
	GetNextQuestion(userID, sessionID int64) (Question, error)
	GetSessionSummary(userID, sessionID int64) (SessionSummary, error)

	GetProfile(userID int64) (Profile, error)
	GetAchievements(userID int64) ([]Achievement, error)

	GetGlobalLeaderboard(limit int) ([]LeaderboardEntry, error)
	GetLeaderboardEntry(userID int64) (LeaderboardEntry, error)
}

type MediaStorage interface {
	GetMedia(bucket, file string) (MediaStream, error)
}

type Service interface {
	Register(email, password, nickname string) error
	Login(email, password string) (string, error)
	UserFromToken(token string) (User, error)

	StartSession(userID int64, mode string) (Session, Question, error)
	SubmitAnswer(userID, sessionID int64, req SubmitAnswer) (AnswerResult, error)
	GetNextQuestion(userID, sessionID int64) (Question, error)
	GetSessionSummary(userID, sessionID int64) (SessionSummary, error)

	GetProfile(userID int64) (Profile, error)
	GetAchievements(userID int64) ([]Achievement, error)

	GetGlobalLeaderboard(limit int) ([]LeaderboardEntry, error)
	GetLeaderboardEntry(userID int64) (LeaderboardEntry, error)

	GetMedia(bucket, file string) (MediaStream, error)
}

type service struct {
	db    DBStorage
	media MediaStorage
}

func New(db DBStorage, media MediaStorage) Service {
	return &service{db: db, media: media}
}

func (s *service) Register(email, password, nickname string) error {
	return s.db.CreateUser(email, password, nickname)
}

func (s *service) Login(email, password string) (string, error) {
	return s.db.LoginUser(email, password)
}

func (s *service) UserFromToken(token string) (User, error) {
	return s.db.GetUserFromToken(token)
}

func (s *service) StartSession(userID int64, mode string) (Session, Question, error) {
	return s.db.StartSession(userID, mode)
}

func (s *service) SubmitAnswer(userID, sessionID int64, req SubmitAnswer) (AnswerResult, error) {
	return s.db.SubmitAnswer(userID, sessionID, req)
}

func (s *service) GetNextQuestion(userID, sessionID int64) (Question, error) {
	return s.db.GetNextQuestion(userID, sessionID)
}

func (s *service) GetSessionSummary(userID, sessionID int64) (SessionSummary, error) {
	return s.db.GetSessionSummary(userID, sessionID)
}

func (s *service) GetProfile(userID int64) (Profile, error) {
	return s.db.GetProfile(userID)
}

func (s *service) GetAchievements(userID int64) ([]Achievement, error) {
	return s.db.GetAchievements(userID)
}

func (s *service) GetGlobalLeaderboard(limit int) ([]LeaderboardEntry, error) {
	return s.db.GetGlobalLeaderboard(limit)
}

func (s *service) GetLeaderboardEntry(userID int64) (LeaderboardEntry, error) {
	return s.db.GetLeaderboardEntry(userID)
}

func (s *service) GetMedia(bucket, file string) (MediaStream, error) {
	return s.media.GetMedia(bucket, file)
}
