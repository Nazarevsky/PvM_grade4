package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"lab4/db"
	"lab4/db/controllers"
	"lab4/db/models"
	"log"
)

const apiKey string = "7101583543:AAGuNv98BUUGELVc3bXm6UKPTHBaTvipDRo"

type Bot struct {
	bot                    *tgbotapi.BotAPI
	disciplineController   controllers.Discipline
	disciplineAtController controllers.DisciplineAttendance
	taskController         controllers.Task
	teacherController      controllers.Teacher
	studentController      controllers.Student
	gradeController        controllers.Grade

	isDisciplineCreation bool
	isDisciplineName     bool
	isMyDiscipline       bool
	isChooseOperation    bool
}

func New() (Bot, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return Bot{}, err
	}

	bot.Debug = true

	dbController, err := db.New()
	if err != nil {
		return Bot{}, err
	}

	return Bot{
		bot:                    bot,
		disciplineController:   controllers.NewDisciplineController(dbController),
		disciplineAtController: controllers.NewDisciplineAttendanceController(dbController),
		teacherController:      controllers.NewTeacherController(dbController),
		taskController:         controllers.NewTaskController(dbController),
	}, nil
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		//b.myDiscipline(update.Message.Chat.ID, update)
		b.newDiscipline(update.Message.Chat.ID, update)
		b.chooseDiscipline(update.Message.Chat.ID, update)

		if update.Message.Text == "меню" {
			b.isDisciplineCreation = false
			b.isDisciplineName = false
			b.isMyDiscipline = false
			b.OpenMenu(update.Message.Chat.ID, "Повертаємось у меню...")
		}

		if update.Message.Text == "restart" {
			//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg.ReplyMarkup = menuKeyboard
			if _, err := b.bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func (b *Bot) newDiscipline(chatId int64, update tgbotapi.Update) {
	if update.Message.Text == "Нова дисципліна" {
		b.isDisciplineCreation = true
		msg := tgbotapi.NewMessage(chatId, "Введіть назву, або \"меню\", щоб повернутись у меню")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		b.isDisciplineName = true

		b.Send(msg)
	} else if b.isDisciplineCreation {
		b.isDisciplineName = false
		b.isDisciplineCreation = false

		var disciplineId uint64
		err := b.disciplineController.Create(models.Discipline{Name: update.Message.Text})
		if err != nil {
			discipline, err := b.disciplineController.GetByName(update.Message.Text)
			if err != nil {
				log.Panic(err)
				return
			}

			disciplineId = discipline.Id
			return
		}

		teacher, err := b.teacherController.GetById(fmt.Sprintf("%d", update.Message.From.ID))
		if teacher.FullName == "" {
			err = b.teacherController.Create(models.NewTeacher(fmt.Sprintf("%d", update.Message.From.ID), update.Message.From.UserName))
			if err != nil {
				log.Panic(err)
				return
			}

			teacherByFullname, err := b.teacherController.GetByFullName(update.Message.From.UserName)
			if err != nil {
				log.Panic(err)
				return
			}

			teacher = teacherByFullname
			return
		}

		err = b.disciplineAtController.Create(models.NewDisciplineAttendance(0, teacher.TelegramId, disciplineId))
		if err != nil {
			log.Panic(err)
			return
		}

		b.OpenMenu(chatId, fmt.Sprintf("Дисципліну '%s' створено", update.Message.Text))
	}
}

func (b *Bot) myDiscipline(chatId int64, update tgbotapi.Update) {
	if update.Message.Text == "Мої дисципліни" {
		b.isDisciplineCreation = false
		disciplineAts, err := b.disciplineAtController.GetByTgId(fmt.Sprintf("%d", update.Message.From.ID))
		if err != nil {
			log.Panic(err)
			return
		}

		var disciplines string
		var names []string
		for i, attendance := range disciplineAts {
			discipline, err := b.disciplineController.GetById(attendance.DisciplineId)
			if err != nil {
				log.Panic(err)
				return
			}

			var class = "Студент"
			_, err = b.teacherController.GetById(attendance.TelegramId)
			if err == nil {
				class = "Викладач"
			}
			disciplines += fmt.Sprintf("%d. %s (%s)\n", i+1, discipline.Name, class)
			names = append(names, discipline.Name)
		}

		if disciplines == "" {
			disciplines = "Ніякі дисципліни не знайдено"
		} else {
			disciplines = "Знайдено наступні дисципліни, в яких ви берете участь: \n" + disciplines
			disciplines += "\nОберіть дисципліну"
		}

		msg := tgbotapi.NewMessage(chatId, disciplines)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = getDisciplineMenu(names)

		b.Send(msg)
		b.isMyDiscipline = false
		b.isChooseOperation = true
	}
}

func (b *Bot) chooseDiscipline(chatId int64, update tgbotapi.Update) {
	if b.isChooseOperation {

		if _, err := b.teacherController.GetById(fmt.Sprintf("%d", update.Message.From.ID)); err == nil {
			println("teacher")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Оберіть операцію")
			msg.ReplyMarkup = getOperationMenu()
			b.Send(msg)
			// teacher
		} else {
			println("student")
			// not teacher
			//discipline, err := b.disciplineController.GetByName(update.Message.Text)
			//if err != nil {
			//	log.Panic(err)
			//	return
			//}
			//
			//tasks, err := b.taskController.GetByDisciplineId(discipline.Id)
			//if err != nil {
			//	log.Panic(err)
			//	return
			//}
		}
		b.isChooseOperation = false
	}
	// if teacher, create task and set grades (menu)
	// if student, just check grades
}

func (b *Bot) Send(msg tgbotapi.MessageConfig) {
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Panic(err)
		return
	}
}

func (b *Bot) OpenMenu(chatId int64, msgText string) {
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ReplyMarkup = menuKeyboard
	b.isDisciplineCreation = false

	b.Send(msg)
}
