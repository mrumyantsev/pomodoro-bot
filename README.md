# Pomodoro-бот

**Pomodoro**&nbsp;&mdash; это менеджер таймеров с интерфейсом чат-бота, который может быть использован для работы по *методу помидора*, либо для установки пользовательских таймеров. После истечения времени у таймера бот отправит соответствующее уведомление пользователю, который его установил.

Функционал:

- возможность установки количества минут;
- возможность установки текста уведомления;
- неограниченное количество таймеров, устанавливаемых пользователями.

# Метод помидора

[**&laquo;Метод помидора&raquo;**](https://ru.wikipedia.org/wiki/%D0%9C%D0%B5%D1%82%D0%BE%D0%B4_%D0%BF%D0%BE%D0%BC%D0%B8%D0%B4%D0%BE%D1%80%D0%B0) (итал. *tunica del pomodoro*)&nbsp;&mdash; техника управления временем, предложенная Франческо Чирилло в конце 1980-х. Методика предполагает увеличение эффективности работы при меньших временных затратах за счет глубокой концентрации и коротких перерывов. В классической технике отрезки времени&nbsp;&mdash; &laquo;помидоры&raquo; длятся полчаса: 25 минут работы и 5 минут отдыха.

![Метод помидора](./pomodoro-timer-technique.jpg "Кухонный таймер-помидор")

# Поддерживаемые клиенты для ботов

* Клиент для [Telegram](https://en.wikipedia.org/wiki/Telegram_(software))

# Применение

Чтобы начать общение с ботом в первый раз, отправьте ему следующее сообщение:

```
/start
```

Чтобы получить справку по командам для бота, отправьте следующее сообщение:

```
/help
```

Чтобы установить стандартный 25-минутный таймер, отправьте следующий текст:

```
set
```

Для установки таймера на выбранное количество минут и дополнительной установки текста уведомления, отправьте текст:

```
set 55 Иди смотреть стрим!!!
```

Для остановки последнего действующего таймера, отправьте:

```
unset
```

Для остановки всех активных таймеров, отправьте текст:

```
unset all
```