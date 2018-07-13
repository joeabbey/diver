#!/bin/bash

echo "This will create a series of users, along with various collections/grants/roles etc."
echo "Sleeping for 5 seconds, press ctrl+c to stop this script"
sleep 5
echo "Beginning Creation"
./diver ucp

if [ $? -eq 0 ]; then
    echo ""
else
    echo "Error running Diver tests"
    exit 1
fi

echo "Creating nine Users"

./diver ucp auth users create  --name billy --password billyPassword --fullname "William Hubert"
./diver ucp auth users create  --name bobby --password bobbyPassword --fullname "Robert Smythe"
./diver ucp auth users create  --name barry --password barryPassword --fullname "Barry Berry"
./diver ucp auth users create  --name nelson --password nelsonPassword --fullname "Nelson Haverford"
./diver ucp auth users create  --name robert --password robertPassword --fullname "Robert Winston"
./diver ucp auth users create  --name EAlderson --password mrRobot123 --fullname "Elliot Anderson"
./diver ucp auth users create  --name DAlderson --password mrsRobot123 --fullname "Darlene Anderson"
./diver ucp auth users create  --name TWellick --password VP4Lyfe! --fullname "Tyrell Wellick"
./diver ucp auth users create  --name PPrice --password CEO4Lyfe! --fullname "Phillip Price"

echo "Creating two Organisations"

./diver ucp auth org create --name ECorp
./diver ucp auth org create --name AllSafe

echo "Adding Teams to Organisations"

./diver ucp auth teams create --team engineering --org AllSafe --description "Allsafe Engineering team"
./diver ucp auth teams create --team executives --org ECorp --description "ECorp Executive team"
./diver ucp auth teams create --team HumanResources --org ECorp --description "ECorp HR team"
./diver ucp auth teams create --team Finance --org ECorp --description "ECorp Finance team"
./diver ucp auth teams create --team Development --org ECorp --description "ECorp dev team"
./diver ucp auth teams create --team Production --org ECorp --description "ECorp production platform team"

echo "Adding Users to Teams"

./diver ucp auth teams adduser --org ECorp --team executives --user EAlderson
./diver ucp auth teams adduser --org ECorp --team executives --user TWellick
./diver ucp auth teams adduser --org ECorp --team executives --user PPrice
./diver ucp auth teams adduser --org ECorp --team Finance --user PPrice

./diver ucp auth teams adduser --org ECorp --team Development --user DAlderson
./diver ucp auth teams adduser --org ECorp --team Development --user billy
./diver ucp auth teams adduser --org ECorp --team Development --user bobby
./diver ucp auth teams adduser --org ECorp --team Development --user barry
./diver ucp auth teams adduser --org ECorp --team HumanResources --user nelson
./diver ucp auth teams adduser --org ECorp --team Production --user robert
./diver ucp auth teams adduser --org AllSafe --team engineering --user EAlderson
./diver ucp auth teams adduser --org AllSafe --team engineering --user DAlderson

echo "Creating Collections"

./diver ucp auth collections create --name ecorp --parent private
./diver ucp auth collections create --name allsafe --parent private

echo "Applying label Constraints"
ECORPID=$(./diver ucp auth collections list | grep ecorp | awk '{print $2}')
ASAFEID=$(./diver ucp auth collections list | grep allsafe | awk '{print $2}')

./diver ucp auth collections set --key security --value none --type engine --id $ECORPID
./diver ucp auth collections set --key hacked --value completely --type node --id $ECORPID
./diver ucp auth collections set --key honeypot --value installed --type node --id $ECORPID

./diver ucp auth collections set --key node.type --value freebsd --type node --id $ASAFEID
./diver ucp auth collections set --key node.gpu --value nvidia --type node --id $ASAFEID

echo "Cloning roles"
./diver ucp auth roles get --id viewonly >> ecorpro.role
./diver ucp auth roles get --id restrictedcontrol >> ecorprw.role

echo "Adding Read Only roles"
./diver ucp auth roles create --rolename ecorpRO --ruleset ./ecorpro.role
rm ./ecorpro.role
echo "Adding Restricted roles"
./diver ucp auth roles create --rolename ecorpRestricted --ruleset ./ecorprw.role
rm ./ecorprw.role

echo "Creating Grants for Teams"
GRANTROLE=$(./diver ucp auth roles list | grep ecorpRestricted | awk '{print $2}')
GRANTUSER=$(./diver ucp auth teams list --org ecorp | grep executives | awk '{print $2}')
./diver ucp auth grants set --collection $ECORPID --subject $GRANTUSER --role $GRANTROLE
GRANTUSER=$(./diver ucp auth teams list --org ecorp | grep Production | awk '{print $2}')
./diver ucp auth grants set --collection $ECORPID --subject $GRANTUSER --role $GRANTROLE
GRANTUSER=$(./diver ucp auth teams list --org ecorp | grep Development | awk '{print $2}')
./diver ucp auth grants set --collection $ECORPID --subject $GRANTUSER --role $GRANTROLE

echo "Creating Grants for Individual Users"
GRANTUSER=$(./diver ucp auth users list | grep ealderson | awk '{print $2}')
./diver ucp auth grants set --collection $ECORPID --subject $GRANTUSER --role $GRANTROLE
./diver ucp auth grants set --collection $ASAFEID --subject $GRANTUSER --role $GRANTROLE
GRANTUSER=$(./diver ucp auth users list | grep twellick | awk '{print $2}')
./diver ucp auth grants set --collection $ECORPID --subject $GRANTUSER --role $GRANTROLE
GRANTUSER=$(./diver ucp auth users list | grep barry | awk '{print $2}')
./diver ucp auth grants set --collection $ECORPID --subject $GRANTUSER --role $GRANTROLE
GRANTUSER=$(./diver ucp auth users list | grep billy | awk '{print $2}')
./diver ucp auth grants set --collection $ECORPID --subject $GRANTUSER --role $GRANTROLE

