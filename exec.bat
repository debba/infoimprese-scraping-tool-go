rem For you fucking Windows User

set /P query="What are you looking for? "
set /P where="Where? "
set /P filename="Filename "
set /P mmode="Mode "

IF NOT DEFINED filename SET "filename=export"
IF NOT DEFINED mmode SET "mmode=with_email"

infoimprese.exe -q %query% -l %where% -m %mmode% -o risultati/%filename%.csv