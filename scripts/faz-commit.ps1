param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$ArgsList
)

legoloop commit @ArgsList
