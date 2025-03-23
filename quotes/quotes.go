package quotes

import (
	"math/rand"
	"time"
)

// Коллекция цитат мудрости
var wisdomQuotes = []string{
	"Знание — сила. (Фрэнсис Бэкон)",
	"Мудрость начинается с удивления. (Сократ)",
	"Настоящая мудрость в том, чтобы знать, что ты ничего не знаешь. (Сократ)",
	"Образование — это оружие, эффект которого зависит от того, кто его держит в руках и на кого оно направлено. (Иосиф Сталин)",
	"Тот, кто хочет видеть результаты своего труда немедленно, должен идти в сапожники. (Альберт Эйнштейн)",
	"Чтобы дойти до цели, надо прежде всего идти. (Оноре де Бальзак)",
	"Если вы не думаете о будущем, у вас его не будет. (Джон Голсуорси)",
	"Сложнее всего начать действовать, все остальное зависит только от упорства. (Амелия Эрхарт)",
	"Успех — это умение двигаться от неудачи к неудаче, не теряя энтузиазма. (Уинстон Черчилль)",
	"Умный человек не делает сам все ошибки — он дает шанс и другим. (Уинстон Черчилль)",
	"Лучше быть хорошим человеком, ругающимся матом, чем тихой, воспитанной тварью. (Фаина Раневская)",
	"Ваше время ограничено, не тратьте его, живя чужой жизнью. (Стив Джобс)",
	"Если вы думаете, что на что-то способны, вы правы; если думаете, что у вас ничего не получится - вы тоже правы. (Генри Форд)",
	"Всегда выбирайте самый трудный путь — на нем вы не встретите конкурентов. (Шарль де Голль)",
	"Сделай сегодня то, что другие не хотят, завтра будешь жить так, как другие не могут. (Джаред Лето)",
}

// GetRandomQuote возвращает случайную цитату из коллекции
func GetRandomQuote() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return wisdomQuotes[r.Intn(len(wisdomQuotes))]
}
