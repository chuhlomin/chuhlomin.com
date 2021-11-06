---
date: 2020-06-15
---

# Projects directory

Different JetBrain's IDEs offer different ways to store projects: IntelliJ IDEA by default will store the project in `~/IdeaProjects`, PyCharm - in `~/PyCharmProjects`, GoLand - in `~/GoLandProjects`, etc. There is a certain logic in the name, and, supposedly, you can easily find the project by remembering in what IDE you opened it in.

Go itself, until recently, only offered to store projects in directories like `~/go/src/github.com/<owner>/<repo>`. This is unusual, but there is a certain beauty: when you clone someone else's project from GitHub, you know exactly where it will be stored.

I've adopted this idea and now store all the projects I work on in directories like `~/Projects/<owner>/<repo>`.  The difference from the scheme proposed by JetBrains is that the division is by the authors of projects, not by the name of the program in which they were opened.

![Projects directory](project-org.png "Projects directory")

If you can't remember the author, you can look it up like this:

```bash
cd ~/Projects

├── be5invis
│   └── Iosevka
```

#advice
