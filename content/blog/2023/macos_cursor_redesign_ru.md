---
date: 2023-11-24
---

# Редизайн курсора в macOS

![MacOS Sonoma Caps Lock indicator](macos_cursor_redesign.png)

Уже много лет я использую [собственную раскладку](https://github.com/chuhlomin/birminitsa)
с русским языком на слое Caps Lock.
(Почему я не использую системное переключение раскладок по Caps Lock или Клавише 🌐?
Да потому что оно работает через раз.)

В MacOS Sonoma зачем-то добавили индикатор текущей раскладки и включённого Caps Lock.
Необъяснимо, почему эта функция не отключается в настройках.

Хорошо, что её еще можно отключить через терминал:

```bash
sudo defaults write \
    /Library/Preferences/FeatureFlags/Domain/UIKit.plist \
    redesigned_text_cursor \
    -dict-add Enabled \
    -bool NO
```

Но и здесь инженеры Apple не сделали так, чтобы эта команда работала сразу — необходимо перезагрузиться.

#advice #macos
