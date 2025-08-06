package manager

// COMMAND HANDLER
const (
	MsgCMDStart = SpiralDelimiter +
		"🌀 SPIRAL CORE — BOOT SEQUENCE INITIATED\n\n" +
		"⚡ SYSTEM STATUS: 🔋 ONLINE\n\n" +
		"🎯 PRIMARY TARGET: YOUR LIMITS\n\n" +
		"🚀 ACTIVE MODULES:\n" +
		"  " + DebtModuleButton + "\n\n" +
		"⏳ IN DEVELOPMENT:\n" +
		"  " + RecipeModuleButton + "\n\n" +
		"  " + GymModuleButton + "\n\n" +
		"  " + TasksModuleButton + "\n\n" +
		"🔊 HEAR IT?\n" +
		"🌀 THE DRILL IS TURNING...\n\n" +
		SpiralDelimiter +
		"🔥 IF YOU BELIEVE YOU CAN — YOU CAN!\n" +
		"💥 BUT IF YOU DON'T BELIEVE —\n" +
		"   THEN JUST START DRILLING\n" +
		"   UNTIL YOU BREAK THROUGH!!!\n" +
		SpiralDelimiter

	MsgCMDHelp = SpiralDelimiter +
		"🌀 SPIRAL COMMAND TRANSMISSION RECEIVED\n\n" +
		"📡 /help — DISPLAY COMBAT MANUAL\n" +
		"🌀 /start — INITIATE SPIRAL CORE\n\n" +
		"💥 ACTIVE DRILL HUBS:\n" +
		"  " + DebtModuleButton + "\n\n" +
		SpiralDelimiter +
		"⏳ IN DEVELOPMENT:\n" +
		"  " + RecipeModuleButton + "\n\n" +
		"  " + GymModuleButton + "\n\n" +
		"  " + TasksModuleButton + "\n\n" +
		"⚠️ WARNING: IN DEVELOPMENT MODULES ARE LOCKED.\n" +
		"💥 THE POWER TO CHANGE YOUR LIFE\n" +
		"IS THE POWER TO DRILL!\n\n" +
		"💢 WHO THE HELL DO YOU THINK YOU ARE?!\n" +
		SpiralDelimiter

	MsgCMDDebt = SpiralDelimiter +
		"🌀 DEBT DRILL HUB — FULL MANUAL 🌀\n\n" +
		"💥 THIS MODULE TRANSFORMS YOUR DEBTS INTO SPIRAL CONTRACTS —\n" +
		"TARGETS FOR YOUR DRILL TO PIERCE THROUGH.\n\n" +
		"🔧 CORE PROTOCOLS:\n\n" +
		"• " + AddDebtButton + " — Forge a new binding agreement with the future\n" +
		"• " + EditDebtButton + " — Modify the terms of an existing contract\n" +
		"• " + PayDebtButton + " — Balance the spiral by returning energy\n" +
		"• " + DeleteDebtButton + " — Annihilate a contract from existence\n" +
		"• " + ListDebtButton + " — Review the history of all active missions\n\n" +
		SpiralDelimiter +
		"⏳ TEMPORAL DRILLING PROTOCOL:\n" +
		"PAST DATES ARE SEALED. ONLY FUTURE DRILLING PERMITTED.\n\n" +
		"• " + SelectDateButton + " — Set temporal coordinates for contract completion\n" +
		"• " + ReInputYearButton + " — Re-drill the year\n" +
		"• " + ReInputMonthButton + " — Re-drill the month\n" +
		"• " + ReInputDayButton + " — Re-drill the day\n" +
		"• " + RedirectDateButton + " — Lock temporal drill and return to parent protocol\n\n" +
		SpiralDelimiter +
		"⚙️ EDIT CONTRACT PROTOCOL:\n\n" +
		"• " + EditDescButton + " — Edit the contract name\n" +
		"• " + EditAmountButton + " — Adjust the spiral power\n" +
		"• " + EditDateButton + " — Reset temporal coordinates\n" +
		"• " + ConfirmEditButton + " — Deploy updated contract\n\n" +
		SpiralDelimiter +
		"🛡️ UNIVERSAL CONTROL PANEL:\n" +
		"Available in all drill sequences:\n\n" +
		"• " + ConfirmButton + " — Confirm and proceed\n" +
		"• " + BackStepButton + " — Retreat to previous step\n" +
		"• " + CancelButton + " — Abort current drill sequence\n" +
		"• " + MainMenuButton + " — Return to Debt Drill Hub\n\n" +
		SpiralDelimiter +
		"🌀 HOW TO DEPLOY:\n" +
		"1. Enter /start — Activate Spiral Core\n" +
		"2. Press '" + MainMenuButtonGeneral + "' — Enter Command Center\n" +
		"3. Select '" + DebtModuleButton + "' — Activate Debt Drill Hub\n" +
		"4. Initiate your first Spiral Contract\n\n" +
		"💥 YOUR DEBTS ARE NOT LIMITS —\n" +
		"THEY ARE TARGETS FOR YOUR DRILL TO PIERCE!\n\n" +
		"🌀 THE SPIRAL NEVER RETREATS — IT ASCENDS.\n" +
		SpiralDelimiter

	MsgCMDRecipe = SpiralDelimiter +
		"🍲 KITCHEN DRILL PROTOCOL - STANDBY 🍲\n\n" +
		"💥 THIS MODULE WILL ALLOW YOU TO:\n\n" +
		"• 🍅 PLAN MEALS & RECIPES\n" +
		"• 🛒 CALCULATE GROCERY LIST\n" +
		"• 🔬 TRACK BZHU (MACROS)\n" +
		"• 🔥 OPTIMIZE CALORIC INTAKE\n\n" +
		"💥 YOUR KITCHEN IS NOT A PANTRY —\n" +
		"IT'S A BATTLEFIELD FOR NUTRITIONAL SUPERIORITY!\n\n" +
		"🚀 COMING SOON — PREPARE FOR LAUNCH!\n" +
		SpiralDelimiter

	MsgCMDGym = SpiralDelimiter +
		"🏋️ HYPERTROPHY DRILL PROTOCOL - STANDBY 🏋️\n\n" +
		"💥 THIS MODULE WILL ALLOW YOU TO:\n\n" +
		"• 💪 LOG WORKOUTS & EXERCISES\n" +
		"• 🔁 TRACK SETS, REPS, WEIGHT\n" +
		"• 📈 MONITOR PROGRESS & GROWTH\n" +
		"• 🚨 DETECT STAGNATION PATTERNS\n\n" +
		"💥 YOUR GYM IS NOT A ROOM —\n" +
		"IT'S A LABORATORY FOR MUSCLE SPIRALIZATION!\n\n" +
		"🚀 COMING SOON — PREPARE FOR GROWTH OVERDRIVE!\n" +
		SpiralDelimiter

	MsgCMDTask = SpiralDelimiter +
		"📅 TASK DRILL PROTOCOL - STANDBY 📅\n\n" +
		"💥 THIS MODULE WILL ALLOW YOU TO:\n\n" +
		"• 🎯 CREATE & MANAGE TASKS\n" +
		"• ⏳ SET DEADLINES & REMINDERS\n" +
		"• 🚀 PRIORITIZE MISSIONS\n" +
		"• ☠️ ANNIHILATE PROCRASTINATION\n\n" +
		"💥 YOUR TO-DO LIST IS NOT A CHORE —\n" +
		"IT'S A SPIRAL TARGET FOR PIERCING!\n\n" +
		"🚀 COMING SOON — PREPARE FOR LIMIT-BREAK!\n" +
		SpiralDelimiter

	MainMenuButtonGeneral = "🌀 DEPLOY COMMAND CENTER 🌀"

	InvalidCommand = SpiralDelimiter +
		"🚨 COMMAND REJECTED BY SPIRAL CORE! 💢\n\n" +
		"💥 SPIRAL CORE ONLY RESPONDS TO PROPER ORDERS!\n" +
		"⚠️ CHECK YOUR INPUT\n\n" +
		"🌀 RETURNING TO SPIRAL CORE...\n" +
		SpiralDelimiter
)

// MAIN MENU
const (
	MsgMainMenu = SpiralDelimiter +
		"🌀 SPIRAL CORE — ONLINE\n\n" +
		"⚡ YOUR LIMITS HAVE BEEN TARGETED AND LOCKED.\n\n" +
		"🌀 DRILL COMMANDER, THE BATTLE BEGINS.\n" +
		"DEPLOY TARGET HUB:\n" +
		SpiralDelimiter

	DebtModuleButton   = "🌀 DEBT DRILL HUB 🌀"
	RecipeModuleButton = "💢 KITCHEN DRILL HUB 💢"
	GymModuleButton    = "💢 GYM DRILL HUB 💢"
	TasksModuleButton  = "💢 TASK DRILL HUB 💢"
)

// GENERAL
const (
	SpiralFormat    = "🌀 %v"
	RageEmoji       = "💢"
	SpiralEmoji     = "🌀"
	SkullEmoji      = "☠️"
	SpiralDelimiter = "─────────🌀─────────\n\n"

	CancelButton   = "🌀✗ DROP DRILLING"
	ConfirmButton  = "🌀↵ LOCK DRILL"
	BackStepButton = "🌀↺ BACK DRILLING"
	MainMenuButton = "🌀 SPIRAL COMMAND CENTER 🌀"

	FailedToCreateKeyboard = SpiralDelimiter +
		"🚨 SPIRAL CONTROL PANEL CRASHED!\n\n" +
		"💥 FAILED TO BUILD BUTTON MATRIX\n\n" +
		"🌀 RETURNING TO SAFE MODE...\n" +
		SpiralDelimiter

	ButtonOnlyMode = SpiralDelimiter +
		"🚨 DRILL COMMAND SYSTEM LOCKED!\n\n" +
		"💥 THIS IS BUTTON-ONLY COMBAT MODE!\n" +
		"⚠️ REAL SPIRAL WARRIORS USE THE DRILL INTERFACE!\n\n" +
		"🌀 ENGAGE TACTICAL COMMAND BUTTONS BELOW!\n" +
		SpiralDelimiter

	InvalidStep = SpiralDelimiter +
		"🚨 ABNORMAL DRILLING SEQUENCE DETECTED!\n\n" +
		"💥 UNKNOWN STEP: %v\n" +
		"⚠️ SPIRAL MATRIX COMPROMISED!\n\n" +
		"🌀 RETURNING TO SAFE MODE...\n" +
		SpiralDelimiter

	InvalidEventType = SpiralDelimiter +
		"🚨 UNIDENTIFIED EVENT SIGNAL!\n\n" +
		"💥 THIS DRILL DOESN'T RECOGNIZE THIS THREAT\n" +
		"⚠️ DEPLOY STANDARD PROTOCOLS ONLY!\n\n" +
		"🌀 RETURNING TO SAFE MODE...\n" +
		SpiralDelimiter

	SessionLost = SpiralDelimiter +
		"🚨 SPIRAL CONNECTION LOST!\n\n" +
		"💥 YOUR DRILL SESSION VANISHED INTO THE VOID\n" +
		"⚠️ REALITY SHIFT DETECTED\n\n" +
		"🌀 RETURNING TO SAFE MODE...\n" +
		SpiralDelimiter

	FailedToGetCallBack = SpiralDelimiter +
		"🚨 CALLBACK SIGNAL LOST IN VOID!\n\n" +
		"💥 BUTTON RESPONSE FAILED TO RETURN\n" +
		"⚠️ THE DRILL IS CONFUSED — COMMUNICATION LOST!\n" +
		"🌀 RETURNING TO SAFE MODE...\n" +
		SpiralDelimiter

	FailedToGetState = SpiralDelimiter +
		"🚨 SPIRAL CONNECTION LOST!\n\n" +
		"💥 FAILED TO LOAD DRILL SESSION\n" +
		"⚠️ SESSION DATA DAMAGED OR LOST\n\n" +
		"🌀 RETURNING TO SAFE MODE...\n" +
		SpiralDelimiter

	MsgFailedToSetSession = SpiralDelimiter +
		"🚨 SPIRAL CONNECTION LOST!\n\n" +
		"💥 FAILED TO SAVE DRILL SESSION\n" +
		"⚠️ SPIRAL ENERGY DISSIPATING\n\n" +
		"🌀 RETURNING TO SAFE MODE...\n" +
		SpiralDelimiter

	HandlerNotFound = SpiralDelimiter +
		"🚨 DRILL PROTOCOL NOT FOUND!\n\n" +
		"💥 NO MATCHING HANDLER FOR EVENT\n" +
		"⚠️ TYPE: %s\n\n" +
		"🌀 PREPARING DEFAULT DRILL SEQUENCE...\n" +
		SpiralDelimiter
)

// DEBT HANDLER
const (
	AddDebtButton    = "🌀 CONTRACT PROTOCOL"
	EditDebtButton   = "🌀 RECALIBRATE PROTOCOL"
	PayDebtButton    = "💥 BALANCE PROTOCOL"
	DeleteDebtButton = "💀 ANNIHILATE PROTOCOL"
	ListDebtButton   = "📜 REVIEW CONTRACT LOG"

	EditDescButton    = "🌀 RE-SET CONTRACT NAME"
	EditAmountButton  = "💥 RE-SET SPIRAL POWER"
	EditDateButton    = "⏳ RE-SET TEMPORAL COORDINATES"
	ConfirmEditButton = "🌀↵ DEPLOY MODIFIED CONTRACT"

	RedirectDebtButton = "🌀↵ LOCK DRILLING TARGET"

	MsgDebtMenu = SpiralDelimiter +
		"🌀 SPIRAL DEBT MODULE — ONLINE\n\n" +
		"💥 YOUR DEBTS ARE NOT LIMITS —\n" +
		"THEY ARE TARGETS FOR YOUR DRILL TO PIERCE!\n\n" +
		"⚡ INITIATE THE BURST:\n\n" +
		"🌀 AWAITING DRILL ORDERS!\n" +
		SpiralDelimiter

	ListDebtFormat = "%s %s\n\t" +
		"💥 SPIRAL POWER: %s₽\n\t" +
		"%s\n\n"

	ListTotalAmountFormat       = "💥 TOTAL SPIRAL POWER REQUIRED: %s₽\n\n"
	ListReturnDateFormat        = "⏳ D-DAY: %d DAYS REMAINING"
	ListReturnDateExpiredFormat = "🚨 ANTI-SPIRAL THREAT (%d DAYS)"
	ReturnDateNil               = "🌌 D-DAY: UNLIMITED BATTLEFIELD"

	MsgAddDescription = "🌀 INITIATE SPIRAL CONTRACT PROTOCOL...\n\n" +
		"💥 INPUT CONTRACT NAME:"

	MsgPayStart = "🌀 INITIATE SPIRAL BALANCE PROTOCOL...\n\n" +
		"💥 SELECT SPIRAL CONTRACT TO BALANCE DRILLING"

	MsgDeleteStart = "🌀 INITIATE SPIRAL ANNIHILATE PROTOCOL...\n\n" +
		"💥 SELECT SPIRAL CONTRACT FOR ERASE DRILLING:"

	MsgEditStart = "🌀 INITIATE SPIRAL RECALIBRATE PROTOCOL...\n\n" +
		"💥 SELECT SPIRAL CONTRACT EDIT DRILLING:"

	MsgAddAmount = "🌀 CONTRACT NAME LOCKED: %s\n\n" +
		"🌀 INITIATE QUANTUM DRILLING!\n\n" +
		"💥 INPUT SPIRAL POWER:"

	MsgAddDate = "🌀 SPIRAL POWER LOCKED: %s₽\n\n"

	MsgEditMenu = "🌀 RECALIBRATE PROTOCOL READY!\n\n" +
		"CHOOSE COMPONENT TO DEEP-DRILLING:"

	MsgEnterAmount = "🌀 INITIATE QUANTUM DRILLING!\n\n" +
		"💥 INPUT NEW SPIRAL POWER:"
	MsgEditAmount = "🌀 CORE DRILLING SUCCESS!\n\n" +
		"💥 SPIRAL POWER RECALIBRATED TO %s₽!\n\n" +
		"🌀 RETURNING TO RECALIBRATE DRILL SEQUENCE..."

	MsgEnterDescription = "🌀 INITIATE CORE DRILLING!\n\n" +
		"💥 INPUT NEW CONTRACT NAME:"
	MsgEditDescription = "🌀 CORE DRILLING SUCCESS!\n\n" +
		"💥 CONTRACT RECALIBRATED TO: %s\n\n" +
		"🌀 RETURNING TO RECALIBRATE DRILL SEQUENCE..."

	MsgEnterDate = "🌀 INITIATE TEMPORAL DRILLING!\n\n" +
		"💥 D-DAY AWAITS YOUR COMMAND!"
	MsgEditDate = "🌀 TEMPORAL DRILLING SUCCESS!\n\n" +
		"💥 D-DAY RECALIBRATED TO: %s\n\n" +
		"🌀 RETURNING TO RECALIBRATE DRILL SEQUENCE..."

	MsgSavedDebt = "🌀 SPIRAL CONTRACT DEPLOYED!\n\n" +
		"🌀 CONTRACT: %s\n" +
		"💥 SPIRAL POWER: %s₽\n" +
		"⏳ D-DAY: %s\n\n" +
		"🌀 THIS CONTRACT IS NOW PART OF THE DRILL LOG\n" +
		"🌀 RETURNING TO COMMAND SEQUENCE..."

	MsgDebtSelected = "🌀 SPIRAL CONTRACT LOCKED!\n\n" +
		"🌀 CONTRACT: %s\n" +
		"💥 SPIRAL POWER: %s₽\n\n" +
		"%s\n\n" +
		"🌀 TARGET LOCKED — NEXT PROTOCOL PENDING\n"

	MsgConfirmDeleteWarning = "☠️ ANNIHILATE DRILL SEQUENCE INITIATED!\n\n" +
		"💀 THIS WILL ERASE THE SPIRAL CONTRACT FROM EXISTENCE\n\n" +
		"🚨 WARNING: THIS ACTION CANNOT BE UNDONE\n" +
		"⚡ THE DRILL WILL PIERCE THROUGH SPACE-TIME\n\n" +
		"🌀 COMMIT TOTAL ANNIHILATION?"

	MsgDeleteDebt = "💀 SPIRAL CONTRACT ERASED!\n\n" +
		"🌀 CONTRACT: %s\n" +
		"💥 SPIRAL POWER: %s₽\n\n" +
		"🌀 THE CONTRACT HAS BEEN DRILLED OUT OF REALITY\n" +
		"🌀 RETURNING TO COMMAND SEQUENCE..."

	MsgPayConfirm = "🌀 BALANCE PROTOCOL FINAL LOCK!\n\n" +
		"🌀 CONTRACT: %s\n\n" +
		"💥 SPIRAL POWER: %s₽\n" +
		"🌀 PAYMENT ENERGY: %s₽\n" +
		"💥 RESIDUAL POWER: %s₽\n\n" +
		"🌀 COMPLETE?"
	MsgPayToDelete = "💥 SPIRAL CONTRACT ANNIHILATED!\n\n" +
		"🌀 CONTRACT:\"%s\" ERASED FROM EXISTENCE\n\n" +
		"🌀 RETURNING TO COMMAND SEQUENCE..."
	MsgPayToUpdate = "🌀 IF THE DEBT IS THIS BIG…\n" +
		"THEN OUR DRILL MUST BE EVEN BIGGER!\n\n" +
		"🌀 CONTRACT: %s\n" +
		"💥 RESIDUAL SPIRAL POWER: %s₽\n\n" +
		"🌀 RETURNING TO COMMAND SEQUENCE..."

	MsgFinishEdit = "🌀 RECALIBRATE PROTOCOL COMPLETE!\n\n" +
		"💥 SPIRAL CONTRACT: %s\n" +
		"🌀 SPIRAL POWER: %s₽\n\n" +
		"%s\n\n" +
		"🌀 RETURNING TO COMMAND SEQUENCE..."

	MsgInvalidDescriptionEmpty = SpiralDelimiter +
		"🚨 SPIRAL CONTRACT REJECTED: NAME FIELD EMPTY!\n\n" +
		"💥 A NAMELESS CONTRACT IS A VOID IN REALITY\n\n" +
		"🌀 RE-ENTER WITH FOCUS:\n" +
		SpiralDelimiter

	MsgInvalidDescriptionLength = SpiralDelimiter +
		"🚨 SPIRAL CONTRACT REJECTED: NAME EXCEEDS LIMIT!\n\n" +
		"💥 CURRENT: %d/1000 CHARACTERS\n" +
		"⚠️ COMPRESS YOUR BATTLE CHANT INTO A TIGHTER SPIRAL\n\n" +
		"🌀 RE-ENTER WITH MAXIMUM DRILLING EFFICIENCY:\n" +
		SpiralDelimiter

	MsgInvalidAmountEmpty = SpiralDelimiter +
		"🚨 SPIRAL ENERGY SIGNATURE INVALID!\n\n" +
		"💥 INPUT LOST IN THE VOID\n\n" +
		"⚠️ INPUT SPIRAL POWER TO IGNITE THE BURST\n\n" +
		"🌀 RE-ENTER VALID SPIRAL POWER:\n" +
		SpiralDelimiter

	MsgInvalidAmountConvertErr = SpiralDelimiter +
		"🚨 SPIRAL ENERGY SIGNATURE INVALID!\n\n" +
		"💥 ONLY POSITIVE NUMBERS CAN PIERCE THE HEAVENS\n\n" +
		"🌀 RE-ENTER VALID SPIRAL POWER:\n" +
		SpiralDelimiter

	MsgToLargeAmount = SpiralDelimiter +
		"🚨 SPIRAL ENERGY SIGNATURE INVALID!\n\n" +
		"💥 SPIRAL POWER EXCEEDS COSMIC LIMIT!\n" +
		"⚠️ MAX: %s₽ — NO MORE, NO LESS\n\n" +
		"🌀 RE-ENTER VALID SPIRAL POWER:\n" +
		SpiralDelimiter

	MsgDateNotSet = SpiralDelimiter +
		"🚨 TEMPORAL COORDINATES LOST!\n\n" +
		"💥 DRILL CANNOT PIERCE THE VOID OF TIME\n\n" +
		"🌀 RETURNING TO CORE PROTOCOL...\n" +
		SpiralDelimiter

	MsgFailedToExtractDebtId = SpiralDelimiter +
		"🚨 SPIRAL SYNC FAILED!\n\n" +
		"💥 CONTRACT IDENTIFIER CORRUPTED!\n" +
		"⚠️ SPIRAL CONTRACT LOCK FAILED — MATRIX INTEGRITY COMPROMISED\n\n" +
		"🌀 REBOOTING DRILL PROTOCOLS...\n" +
		SpiralDelimiter

	MsgFailedToGetDebt = SpiralDelimiter +
		"🚨 SPIRAL SYNC FAILED!\n\n" +
		"💥 SPIRAL CONTRACT NOT FOUND IN MATRIX!\n\n" +
		"⚠️ THE TARGET HAS VANISHED FROM REALITY\n\n" +
		"🌀 REBOOTING DRILL PROTOCOLS...\n" +
		SpiralDelimiter

	MsgUserIdNotEqualDebtId = SpiralDelimiter +
		"🚨 DRILL COLLISION DETECTED!\n\n" +
		"💥 THIS SPIRAL CONTRACT BELONGS TO ANOTHER PILOT\n\n" +
		"⚠️ YOU CANNOT PIERCE ANOTHER MAN'S SOUL ⚔️\n\n" +
		"🌀 REBOOTING DRILL PROTOCOLS...\n" +
		SpiralDelimiter

	MsgFailedToSaveDebt = SpiralDelimiter +
		"🚨 SPIRAL CONTRACT REGISTRY REJECTED!\n\n" +
		"💥 DRILLING FAILURE!\n\n" +
		"⚠️ SPIRAL COLLAPSE DETECTED — UNIVERSE RESISTS OUR DRILL\n" +
		"🌀 REBOOTING DRILL PROTOCOLS...\n" +
		SpiralDelimiter

	MsgFailedToDeleteDebt = SpiralDelimiter +
		"🚨 SPIRAL CONTRACT ANNIHILATION REJECTED!\n\n" +
		"💥 DRILLING FAILURE!\n\n" +
		"⚠️ SPIRAL COLLAPSE DETECTED — UNIVERSE RESISTS OUR DRILL\n" +
		"🌀 REBOOTING DRILL PROTOCOLS...\n" +
		SpiralDelimiter

	MsgFailedToUpdateDebt = SpiralDelimiter +
		"🚨 SPIRAL CONTRACT RE-DEPLOYMENT REJECTED!\n\n" +
		"💥 DRILLING FAILURE!\n\n" +
		"⚠️ SPIRAL COLLAPSE DETECTED — UNIVERSE RESISTS OUR DRILL\n" +
		"🌀 REBOOTING DRILL PROTOCOLS...\n" +
		SpiralDelimiter

	FailedToGetDebts = SpiralDelimiter +
		"🚨 SPIRAL CORE CORRUPTED!\n\n" +
		"💥 SPIRAL MATRIX OFFLINE!\n" +
		"⚠️ FAILED TO SCAN CONTRACTS!\n\n" +
		"🌀 REBOOTING DRILL PROTOCOLS...\n" +
		SpiralDelimiter
)

var (
	DebtTitles = []string{
		"💢 YOUR MONEY OWES YOU MONEY?!",
		"💢 DEBTS TRIED TO HIDE... BUT WE DUG TOO DEEP!",
		"💢 INTEREST RATES? MORE LIKE INTEREST ENEMIES!",
		"💢 CREDITORS DREAM OF YOUR FAILURE - DRILL THEIR DREAMS!",
		"💢 THESE NUMBERS DARE CALL THEMSELVES NEGATIVE?!",
		"💢 FINANCIAL GRAVITY WON'T HOLD US DOWN!",
		"💢 BANKERS HATE THIS ONE SIMPLE DRILL!",
		"💢 YOUR BALANCE SHEET IS ABOUT TO GET BALANCED... WITH A DRILL!",
		"💢 DEBT-TO-INCOME RATIO? MORE LIKE DEBT-TO-OBLIVION!",
		"💢 COMPOUND INTEREST MEETS COMPOUND DRILLING!",
	}

	MotivationalPhrases = []string{
		"💢 PAYMENTS ARE JUST WEAK FUTURE VERSIONS OF YOU!",
		"💢 YOUR CREDIT SCORE IS ABOUT TO GET EPIC!",
		"💢 THEY SAID 'MINIMUM PAYMENT' - WE SAID 'MAXIMUM DRILL'!",
		"💢 COLLECTION AGENTS FEAR YOUR DRILL!",
		"💢 APR STANDS FOR 'ABOUT TO GET PIERCED, REBELS!'",
		"💢 FORECLOSURE? MORE LIKE FORE-DRILL-SURE!",
		"💢 YOUR BANK STATEMENT NEVER SAW THIS COMING!",
		"💢 INTEREST ACCRUED? TIME TO ACCRUE SOME DRILLING!",
		"💢 PRINCIPAL BALANCE MEETS PRINCIPAL DRILL!",
		"💢 LIABILITIES ARE JUST ASSETS WAITING TO BE DRILLED!",
	}

	NoDebtsPhrases = []string{
		"YOUR DEBTS HAVE BEEN DRILLED INTO COSMIC DUST!",
		"CREDITORS WEEP AS YOUR BALANCE SHINES!",
		"THE WORD 'INTEREST' NOW INTERESTS NO ONE!",
		"FINANCIAL FREEDOM ACHIEVED - LIKE A BOSS!",
		"YOUR MONEY FINALLY WORKS FOR YOU (FOR ONCE)!",
		"DEBT: [NULL]. DRILL: [MAXIMUM]. LIFE: [EPIC].",
		"BANKS: CONFUSED. YOU: VICTORIOUS.",
		"LOAN OFFICERS NOD IN RELUCTANT RESPECT!",
		"THE 'CREDIT' IN 'CREDIT SCORE' NOW CREDITS YOU!",
		"COMPOUNDING RETURNS > COMPOUNDING INTEREST!",
	}
)

// DATE HANDLER
const (
	MsgStartDateFlow = "🌀 ACTIVATE TEMPORAL DRILL SEQUENCE!\n\n" +
		"⏳ TEMPORAL DRILL CHARGE: 100%\n\n" +
		"🌀 D-DAY AWAITS YOUR COMMAND!"

	MsgSetYear = "🌀 YEAR DRILL ENGAGED!\n\n" +
		"⏳ SELECT DESTINATION YEAR"

	MsgSetMonth = "🌀 YEAR %d LOCKED!\n\n" +
		"⏳ SELECT DESTINATION MONTH"

	MsgSetDay = "🌀 MONTH %s LOCKED!\n\n" +
		"⏳ SELECT DESTINATION DAY"

	MsgEmptyDay = SpiralDelimiter +
		"🚨 TIME PARADOX DETECTED!\n\n" +
		"💥 SPECIFIED MONTH DOESN'T EXIST\n" +
		"⚠️ YOUR DRILL PIERCED A TIME HOLE\n" +
		"🌀 RE-INITIATE DATE INPUT FROM START!\n" +
		SpiralDelimiter

	MsgRedirect = "⏳ TEMPORAL COORDINATES LOCKED!\n\n" +
		"🌀 YEAR: %d\n" +
		"⏳ MONTH: %s\n" +
		"🌀 DAY: %d\n\n" +
		"⏳ TEMPORAL-DRILL PRIMED FOR DEPLOYMENT!"

	MsgFailedRedirection = SpiralDelimiter +
		"🚨 THE SPIRAL HAS NO DIRECTION!\n\n" +
		"💥 IT DOES NOT ASCEND\n" +
		"💥 IT DOES NOT RETREAT\n\n" +
		"⚠️ THIS DRILL SEQUENCE EXISTS IN VOID\n\n" +
		"🌀 ABORTING — RETURNING TO SAFE REALITY...\n" +
		SpiralDelimiter

	MsgInvalidYear = SpiralDelimiter +
		"🚨 TEMPORAL ANOMALY DETECTED!\n\n" +
		"💥 INVALID YEAR FORMAT!\n" +
		"🌀 RE-DRILLING YEAR!\n" +
		SpiralDelimiter

	MsgInvalidMonth = SpiralDelimiter +
		"🚨 TEMPORAL ANOMALY DETECTED!\n\n" +
		"💥 INVALID MONTH FORMAT!\n" +
		"🌀 RE-DRILLING MONTH!\n" +
		SpiralDelimiter

	MsgInvalidDay = SpiralDelimiter +
		"🚨 TEMPORAL ANOMALY DETECTED!\n\n" +
		"💥 INVALID DAY FORMAT!\n" +
		"🌀 RE-DRILLING DATE!\n" +
		SpiralDelimiter

	MsgDateInPast = SpiralDelimiter +
		"🚨 TEMPORAL ANOMALY DETECTED!\n\n" +
		"💥 CURRENT DRILL TIMESTAMP: %s\n\n" +
		"💥 ATTEMPTED PAST-DRILLING: %s\n\n" +
		"⚠️ TEMPORAL PROTOCOL VIOLATION: PAST DRILLING FORBIDDEN\n\n" +
		"🌀 ONLY FUTURE DRILLING PERMITTED\n\n" +
		"🌀 RESTART TEMPORAL DRILL PROTOCOL...\n" +
		SpiralDelimiter

	SelectDateButton = "⏳ SET D-DAY"

	ReInputYearButton  = "🌀↻ RE-DRILL YEAR"
	ReInputMonthButton = "🌀↻ RE-DRILL MONTH"
	ReInputDayButton   = "🌀↻ RE-DRILL DAY"

	RedirectDateButton = "🌀↵ LOCK TEMPORAL DRILL" // REDIRECT TO PARENT HANDLER
)
