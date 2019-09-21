rem For you fucking Windows User

set /P query="What are you looking for? "
set /P where="Where? "
set /P filename="Filename "
set /P mode="Mode "

IF NOT DEFINED sleep SET "sleep=1"
IF NOT DEFINED filename SET "filename=1"
IF NOT DEFINED mode SET "mode=search_by_name"

./infoimprese -q %query% -l %where% -m %mode% -o risultati/%filename%.csv