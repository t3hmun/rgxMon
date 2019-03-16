# rgxmon

Console app that monitors a regex you are editing and spams the matches *with groups* to the console.
All matches are abbreviated (middle removed) to fit the width of the console, should work on most Go compatible syste,s.

```bash
rgxmon regexfile targetfile
```

The `regexfile` is a text file with the regex you are authoring.

The `targetfile` is the file the regex is being tested against.


## Rationale

I often write long regex with groups that capture specific fields.
This lets me see the result as I edit each group to match the bit I want.

Also this is my first Go program apart from 'hello world'.


## Bugs/Quirks

The terminal width faield on my Windows Quake-style Bash ConEmu.
I've set it to default to 80 chars width whern it fails.


## License

The [Wiccan Rede](https://en.wikipedia.org/wiki/Wiccan_Rede) about sums it up:

>An’ ye harm none, do what ye will

MIT license for those of you need a real licence.
