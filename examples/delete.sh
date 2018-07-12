#!/bin/bash

echo "This will cleanup all of the users/collections etc.. created by the create.sh script"
echo "Sleeping for 5 seconds, press ctrl+c to stop this script"
sleep 5
echo "Beginning Cleanup"
./diver ucp

if [ $? -eq 0 ]; then
    echo ""
else
    echo "Error running Diver tests"
    exit 1
fi

echo "Removing nine Users"

./diver ucp auth users delete  --name billy 
./diver ucp auth users delete  --name bobby 
./diver ucp auth users delete  --name barry 
./diver ucp auth users delete  --name nelson 
./diver ucp auth users delete  --name robert 
./diver ucp auth users delete  --name EAlderson 
./diver ucp auth users delete  --name DAlderson 
./diver ucp auth users delete  --name TWellick 
./diver ucp auth users delete  --name PPrice

echo "Removing two Organisations"

./diver ucp auth org delete --name ECorp
./diver ucp auth org delete --name AllSafe

echo "Removing Collections"

ECORPID=$(./diver ucp auth collections list | grep ecorp | awk '{print $2}')
ASAFEID=$(./diver ucp auth collections list | grep allsafe | awk '{print $2}')
./diver ucp auth collections delete --id $ECORPID
./diver ucp auth collections delete --id $ASAFEID

echo "Removing Roles"
ECORPROLERO=$(./diver ucp auth roles list | grep ecorpRO | awk '{ print $2 }')
ECORPROLERW=$(./diver ucp auth roles list | grep ecorpRestricted | awk '{ print $2 }')

./diver ucp auth roles delete --id $ECORPROLERO
./diver ucp auth roles delete --id $ECORPROLERW

