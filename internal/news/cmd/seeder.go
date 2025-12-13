package main

import (
	"context"
	"fmt"

	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const timeFormat string = "2006-01-02"

var args = []db.CreateNewsWithPublishDateParams{
	{
		CreatorUsername: "admin",
		Title:           `Snack Kingdom Unveils Mission Statement`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-12", timeFormat), Valid: true},
		Content:         `SNACK KINGDOM — In a surprising move toward corporate flair, SNK has announced its official mission statement, distilled directly from member feedback: “SUPER NUTS KRAZY.” Leadership described the phrase as deeply aspirational and legally meaningless, but emotionally accurate. Members praised the branding for capturing Snack Kingdom’s identity — unpredictable, joyful, and slightly dangerous around cookies. No actual rules were harmed in its creation. —Catty von Catacean`,
	},

	{
		CreatorUsername: "admin",
		Title:           `North King’s Greatest Achievement Was Waiting`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-12", timeFormat), Valid: true},
		Content:         `SNACK KINGDOM — When asked about his proudest moment, North King didn’t name a battle or a burn. He named restraint. For nearly two months, he survived on M&Ms alone, saving every item for SVS preparation. When the moment came, everything was used at once, efficiently and without waste. He is proud. So are we. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `HAVOC Holiday Identity Crisis — Whaleslaughter Accused of Being Too Cute and Sweet`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-12", timeFormat), Valid: true},
		Content:         `The state was left emotionally unbalanced today after a player dared to claim Whaleslaughter’s Christmas songs were "too cute and sweet." Witnesses reportedly dropped candy canes and panic-streamed in disbelief. Sources confirm Whaleslaughter responded by sharpening a candy cane while muttering about snow and fire. Confusion deepened when a haunting Italian track surfaced — described as "weirdly beautiful" and extremely suspicious. Curious citizens may review the evidence in this song. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Snack Kingdom Adotta una Politica Anti-Snack Mollicci`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-12", timeFormat), Valid: true},
		Content:         `SNACK KINGDOM — In un’intervista condotta in italiano impeccabile, Matti (Popo Arrabiato) ha chiarito la filosofia di SNK: benevolenza assoluta, ma tolleranza zero per snack mollicci. Chi porta l’uvetta non viene punito, bensì indirizzato verso il celebre “Mix di Riflessione.” Secondo Matti, ordine e croccantezza possono coesistere. Catty prende nota, con rispetto e una certa fame. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Snack Kingdom Leaders Reflect on Pride, Crunch, and Staying Power`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-12", timeFormat), Valid: true},
		Content:         `In a rare sit-down interview, Snack Kingdom leaders were asked about their proudest moments. ADM pointed not to recent drama, but to the long-ago moment the state chose peace. Nika spoke of endurance — SNK still standing, still family. Almas praised unity itself: bringing people together around one shared goal, preferably edible. Pride, it seems, comes in many flavors. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Catty Interviews Core, Receives Candy-Colored Answers`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-12", timeFormat), Valid: true},
		Content:         `During a sit-down interview, Catty asked Core a series of questions about Snack Kingdom and state pride. Core responded with the following statements, in order: “🤔🍬”, “❓🍭”, and “🤨🍫”. When asked to clarify, Core added, “💭🍪”. Catty reports the answers were vibrant, consistent, and impossible to misquote. Interpretation has been deferred. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Whaleslaughter Briefly Puzzled by Ongoing Tension`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-11", timeFormat), Valid: true},
		Content:         `Supreme President Whaleslaughter was observed today pausing mid-governance, mildly confused by the intensity of certain reactions. Rather than investigate further, she reportedly stepped into a brief emotional funk — the rhythmic kind — and let it resolve itself somewhere between a bassline and a chorus. A song now exists. 🎵 Why Can’t We Be Frenemies No policies were changed. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `The Acrobatics of Accountability`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-11", timeFormat), Valid: true},
		Content:         `A new trend swept through state diplomacy today: the I-did-nothing-wrong-but-don’t-be-upset apology. It’s a charming routine in which the speaker insists their words were harmless, impersonal, and entirely justified — yet somehow expects the reassurance to land softly. Catty applauds the flexibility required. Not everyone can twist themselves into a pretzel while claiming they never bent at all. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Unofficial Policy Leaks From Presidential Aide`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-11", timeFormat), Valid: true},
		Content:         `An aide was overheard today giving the most concise summary yet of Whaleslaughter’s current stance toward leader chat: “Real business goes to the president. Opinions can form a line at the bathroom.” Catty sees no reason to confirm the quote. When a remark matches the mood of the state this perfectly, its accuracy becomes a secondary concern. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Transfer Period Ends; Noise Levels in Leader Chat Reach Historic Highs`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-11", timeFormat), Valid: true},
		Content:         `The transfer window closed smoothly today, even as Nut Leader chat delivered a dramatic reading of its grievances. Amid the noise, Sandwich and ADM stood out with rare reason and actual information. Whaleslaughter exited once the discussion devolved into finger-pointing, leaving cooler heads to steady the state. Only tempers and egos burned. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Cold Shoulder Coalition Issues Frostbite Warning`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-11", timeFormat), Valid: true},
		Content:         `A new climate advisory was issued today after the Cold Shoulder Coalition released yet another blacklist update. Analysts noted a curious pattern: several names were added despite no record of burns, battles, or even basic interaction. Experts classify this as 'ambient frostbite' — a condition in which the temperature drops because the Coalition enjoys generating cold fronts on principle. No fire is required; a mild breeze of imagined offense is often enough. Citizens are encouraged to dress warmly and avoid prolonged exposure to the draft. Those caught downwind may experience numbness, irritation, or sudden urges to roll their eyes. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Jane: Power, Poise, and Deeply Questionable Cookie Taste`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-08", timeFormat), Valid: true},
		Content:         `SNACK KINGDOM — Jane is the kind of leader who builds alliances, steadies presidents, and makes a chaotic state feel survivable. A former president herself, she works relentlessly behind the scenes, networking, encouraging, and quietly strengthening everything she touches. This is also, unfortunately, the same woman who eats raisins in cookies. Great leaders are not flawless. They’re just excellent enough that we forgive them anyway. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Monkey and TSK KING 04 Close SVS with Handshake`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-06", timeFormat), Valid: true},
		Content:         `SVS ended on a rare but welcome sight as Monkey and TSK KING 04 met at center field for a handshake after the final scores were locked. The moment marked a clean close to a hard-fought event, with both players acknowledging effort, timing, and execution on both sides. No speeches. No scoring theatrics. Just a nod, a shake, and mutual respect. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Monkey Wins SVS Without Pet Buffs, Refuses to Explain Further`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-06", timeFormat), Valid: true},
		Content:         `In what experts are calling a personal attack on effort itself, Monkey reportedly secured the SVS win without using a single pet buff. Sources confirm the pets were present, loyal, and emotionally prepared. They were simply not selected. Analysts are now divided over whether this was strategic restraint or psychological warfare. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `2128 Crushes 2184, HAVOC R4s Battle for Right to Not Be Supreme President`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-06", timeFormat), Valid: true},
		Content:         `With 2128 now ahead of State 2184 by an irresponsible number, the alliance has secured the right to appoint the next Supreme President. Unfortunately, this has triggered internal diplomatic disaster as multiple R4s are currently fighting to avoid being crowned. Leaked chat reports Raistlin declaring, "Not a chance in hell." Meanwhile, insiders confirm Whaleslaughter and A Sandwich are insisting they absolutely do not want the role while both just happen to be standing extremely close to the throne at all times.—Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Monkey Snatches SVS Crown with 2M Lead, Death Files Emotional Appeal`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-05", timeFormat), Valid: true},
		Content:         `In a finish so tight it caused several keyboards to be smashed in slow motion, Monkey claimed overall victory on Day 5 of SVS Prep, edging out Death by a brutal 2 million points. Witnesses report Baby Death stared at the leaderboard for a full thirty seconds before whispering, "That’s illegal," and gently laying down. Monkey immediately celebrated by existing loudly and holding the castle for the entire 2 hours. —Catty von Catacean, reporting from the scene of the crime`,
	},
	{
		CreatorUsername: "admin",
		Title:           `WOS Discord Did Not Think Whaleslaughter’s Feet Pics Were as Funny as She Did`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-04", timeFormat), Valid: true},
		Content:         `Mods acted within seconds to hide Whaleslaughter’s pretty pedicure after it appeared briefly on the WOS Discord and was immediately judged unacceptable. Sources confirm the removal was swift, quiet, and devastating to morale. The small cartoon bomb nearby is still under internal review, though early analysis suggests it was a visual joke and not a credible threat. Whaleslaughter has not issued a statement, but her toenails remain glossy and undefeated. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Correction Issued: Monkey Was Not Second`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-04", timeFormat), Valid: true},
		Content:         `An official correction has been issued after further review confirmed that Monkey was, in fact, **nowhere near second place** during SVS Prep. The correct standings are as follows: **Angry ???G (SNK)** claimed second place, **Death** secured third, and Monkey bravely explored the lower leaderboard biome at approximately **22nd**. We regret the statistical error. We do not regret the narrative arc. —Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Bean There, Won That: Mr One Stuns SVS Prep, Topples Monkey & Death`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-04", timeFormat), Valid: true},
		Content:         `In a result nobody had penciled in and everyone is now pretending they saw coming, **Mr One** stormed SVS Prep and knocked both Monkey and Death clean off their pedestals. Hailing from **State 2184** and flying the banner of **GÖKTÜRK**, he arrived quietly, stocked with beans, and left carrying the entire event. Analysts are still reviewing the footage frame by frame to determine whether this was strategy, sorcery, or simply Bean Energy at industrial scale. One thing is certain: Monkey and Death both got… beaned. —Catty von Catacean, reporting from a suddenly very quiet leaderboard`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Baby Death Naps Through Day Three as Monkey Steals the Lead`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-03", timeFormat), Valid: true},
		Content:         `In a shocking upset no one actually saw coming, Monkey has officially dethroned Death on day three of SVS prep. Witnesses report Baby Death was last seen clutching a suspiciously sticky lollipop and mumbling something about “just resting his eyes” before vanishing from the leaderboard entirely. Alliance analysts are now split on what happened: • exhaustion-induced nap • sugar crash • or an elaborate fake-out to guilt Monkey later. —Catty von Catacean, reporting live from the nursery crypt`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Death Still Refusing to Be Normal (Again)`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-12-02", timeFormat), Valid: true},
		Content:         `Death woke up today and chose first place again. Two days straight. Monkey made a run for it yesterday and Death simply said, 'Thanks for keeping my spot warm.' The leaderboard briefly believed in miracles and then immediately apologized. At this point we’re considering bolting the crown directly to Death’s furnace. Monkey remains in second, watching… learning… and definitely plotting something unholy. —Catty von Catacean, reporting live from the crushed dreams department`,
	},
	{
		CreatorUsername: "admin",
		Title:           `State 2128 Enters Collective Food Coma After Legendary Feast`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-29", timeFormat), Valid: true},
		Content:         `State 2128 has officially collapsed into a synchronized food coma following a feast so powerful it neutralized all hostilities and most motor skills. Fighters, farmers, and alliance leaders alike were last seen flattened across snowbanks and tavern floors with belts undone and zero regrets. Reports confirm no battles occurred because everyone was too full to stand, argue, or care. Authority has temporarily transferred to whoever can still reach the dessert table.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `HAVOC Holidays Goes Global — Whaleslaughter Dodges the Spotify Wall`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-28", timeFormat), Valid: true},
		Content:         `After discovering that the Spotify link ghosted half the planet, HAVOC Holidays has officially gone live on Suno, courtesy of Whaleslaughter. Same wild carols and festive menace — now actually playable worldwide. If Suno isn’t your platform of choice, just search Whaleslaughter on your favorite music app and find the chaos wherever you listen. No fake cheer. Just havoc.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Great Dane Drops a Proper Christmas Song — Catty Says 'Finally!'`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-25", timeFormat), Valid: true},
		Content:         `Holiday panic broke out across State 2128 today as Great Dane released what many are calling 'an actual Christmas song' — the warmly traditional In the Light of Christmas. The pure, classic sincerity of the track was reportedly so wholesome that Whaleslaughter audibly gagged at its happy cheer. Catty von Catacean, long-suffering cultural critic, issued a grateful statement: “If you want something that actually *sounds* like Christmas, listen to this. Whaleslaughter’s songs are… well… a different kind of Christmas.”`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Leader Chat Argument Becomes Primetime Drama; President Q Zar Jails Two`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-23", timeFormat), Valid: true},
		Content:         `What began as a routine disagreement in Leader Chat escalated into a full-length soap opera, prompting President Q Zar to throw both Vampire and Artemis into ice jail. Far from distressed, spectators reportedly read every line like it was reality television. Vampire spent his incarceration loudly critiquing presidential policy and offering alternative governance strategies, while Artemis simply waited for the credits to roll.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `🎄 HAVOC HOLIDAYS FIRST REVIEW — Album 'A Vibe,' Santa’s Snacks Dragged Publicly`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-23", timeFormat), Valid: true},
		Content:         `While doing household chores and performing the sacred 'cleaning test,' a real world listener delivered the first official review of *HAVOC Holidays.* She praised the album’s variety and overall vibe but singled out 'Santa’s Snacks' as the auditory equivalent of finding glitter stuck to your carpet forever. Listen to the track here. HAVOC spokespeople insist the song is 'a culinary narrative triumph' and not, in fact, a crime.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Mordukai Responds to 'More Dookie' Nickname Allegations`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-23", timeFormat), Valid: true},
		Content:         `After reports circulated that Jinx had referred to him as “more dookie” in her sleep, Mordukai was asked for comment. He paused, shrugged once, and replied, “Don’t matter much.” Witnesses say he then resumed whatever he was doing without further acknowledgement. —reported by Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `3LGuapo’s Conservative March Comes With Unapologetic 🍆 Energy`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `3LGuapo continues to operate with a slow, calculated approach, crediting Cherry Blossom as the NLO alliance’s true calming force. When asked to summarize his vibe, he responded only with 🍆, leaving strategists to draw their own conclusions. He referred to outside groups simply as 'bullies,' offering no elaboration. —Catty von Catacean, sipping tea and taking notes`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Whaleslaughter Charges Ecclesiastical Turkey Federation With Only Two Emotional Support Sheepdogs`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `Witnesses report that Whaleslaughter initiated direct engagement against forces of the Ecclesiastical Turkey Federation while fielding a roster consisting exclusively of two sheepdogs and an alarming amount of confidence. When questioned, she clarified the maneuver was the result of a dare and therefore 'legally exempt from criticism.' Despite catastrophic troop losses, alliance morale reportedly increased. —reported by Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Queen Zar Upholds Fair Play, But Territory Remains Sacred`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `Queen Zar describes her strategy as fair and names 3LGuapo as the steadying presence within the NLO group. She avoids racing Phil or Salsa unless they stray into her territory while gathering. When discussing other coalitions, she offered a single diplomatic word: 'Peace,' delivered without sarcasm. —reported by Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Mordukai Deals Quiet Retribution, Unmoved by Outside Forces`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `Mordukai plays with restrained timing until provoked, at which point consequences arrive quietly and decisively. His movements intersect with Sin's often—sometimes in alignment, sometimes orthogonal—always intentional. When pressed for thoughts on nearby factions, he shrugged the topic off entirely: 'Don't matter much.' —Catty von Catacean, menace of the written word`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Faceless Drifts Through Races Unbothered, Calls Outside Groups 'Querulous'`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `Faceless maintains a relaxed pace, naming Sin as the calmest member of NLO while embracing a steady 😑 energy. They declined to name opponents they'd avoid in a race, citing insufficient reason to care. When referring to other banners, Faceless chose a single word: 'Querulous,' sparking debate on whether this was shade or science. —Catty von Catacean, sipping tea and taking notes`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Jinx Reportedly Mumbles “More Dookie” in Her Sleep`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `Multiple sources report hearing Jinx mutter the phrase “more dookie” in her sleep, followed by unexplained villain-laughter in the early hours of the morning. The nickname is believed to refer to Mordukai, though Jinx declined to clarify when awake, stating only, “If he wants a new name, he can earn it.” Mordukai was not immediately available for comment. —reported by Catty von Catacean`,
	},
	{
		CreatorUsername: "admin",
		Title:           `BREAKING — BIA Fun Levels Spike, Salt Reserves Critically Low`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `In a shocking turn of events, BIA participants reported having fun—yes, FUN—despite early predictions of full-scale meltdowns. Salt output dropped 73% compared to the previous session, though experts warn levels could surge the moment someone loses 490,000 troops to Dii of SNK again. Analysts are still investigating how 'strategy' snuck into the event without proper clearance.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `HAVOC Holidays Hits Spotify — Seasonal Threats Now Streaming`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `Whaleslaughter is excited to announce Havoc Holidays is now streaming on Spotify as a playlist of festive war anthems forged in late-night battles, icy apocalypses, and questionable group chats. Three tracks are still stuck in copyright purgatory, so the full album will arrive once the paperwork stops crying. Listen and find yourself if you dare: HAVOC Holidays — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Grim Reaper Caught Harvesting HAVOC for BIA — HAVOC Refuses to Let Go`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-22", timeFormat), Valid: true},
		Content:         `Sources confirm the Grim Reaper has been quietly harvesting HAVOC troops for BIA, slipping between forts and collecting souls like Pokémon cards. The Employer That Fumbled and lost him insists he’s on temporary assignment; HAVOC has issued a counter-report: his mail has already been forwarded. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Jinx Returns, Insults Three People and a Mod — Whales Drops Bombshell Era on Sight`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-20", timeFormat), Valid: true},
		Content:         `Jinx vanished for eight months without warning, then returned like nothing happened and immediately offended three people and a moderator. Whales, just seven days into her so-called bombshell era, abandoned self-improvement the moment Jinx sat down. No journaling. No manifestation. She lit a cigar and asked, “Okay, who do we take out first?” — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `BREAKING — Whalerilyn Demands Frost Stars as Proof of Identity`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-17", timeFormat), Valid: true},
		Content:         `After being accused of catfishing by a bewildered visitor from State 2336, Whalerilyn Monroe issued the only reasonable solution: donate frost stars to confirm she’s real. Observers report she delivered this with a smirk sharp enough to cut ice. The poor traveler, buried under question marks and existential dread, has yet to respond. HAVOC economists warn the frost-star market may destabilize. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `BREAKING — Granny Beelzebub Goes Full Rambo`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-16", timeFormat), Valid: true},
		Content:         `With Whaleslaughter temporarily convinced she’s a sweet and sexy bombshell, Granny Beelzebub has taken over as HAVOC’s frontline muscle. Witnesses describe her as ‘pure sinew and divine fury,’ storming through snowdrifts with the energy of a retired general who finally snapped. Asked about her new role, she growled, ‘Somebody’s gotta keep order.’ HAVOC leadership remains in turmoil. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `NLO’s Cherry Blossom Breaks Her Silence`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-16", timeFormat), Valid: true},
		Content:         `In a rare interview with Catty von Catacean, NLO’s Cherry Blossom finally spoke—and instantly reminded the state why the quiet ones are the most dangerous. Asked for her playstyle, she answered with the single, icy word: “Chill.” She named Sami as the calmest in NLO, confirming that serenity may be the alliance’s true weapon. And when asked for an emoji to sum up her vibe, she sent a kissy face, drifting across the frost like a gentle warning.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `BREAKING — ‘Single Jane’ Declares a New Era of Independence`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-16", timeFormat), Valid: true},
		Content:         `In a bold act of self-branding, Jane has officially renamed herself ‘Single Jane,’ prompting immediate speculation, celebration, and a few mid-battle wardrobe adjustments. Observers report the name tag change caused several commanders to blush and walk into snowbanks. Asked for comment, Single Jane simply nodded with regal confidence. The editor recommends hydration and emotional stability for all players impacted by this development.— Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `BREAKING — Salsa Activates Anti-Recon After Reconning Everyone`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-16", timeFormat), Valid: true},
		Content:         `After splattering himself across HAVOC last night like an overexcited enchilada, local condiment Salsa has slammed on anti-recon, declaring himself a spicy mystery. Experts call it ‘peak Taco-Bell-energy,’ comparing it to wearing sunglasses indoors and thinking it’s stealth. Salsa remains at large, last seen making a suspicious run for the border - HAVOC is not chasing him.— Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `PAID ADVERTISEMENT — Sun God Nika Launches Popcorn Empire`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-14", timeFormat), Valid: true},
		Content:         `Tired of watching state drama on an empty stomach? Sun God Nika has descended from the heavens with a solution: Nika’s Divine Popcorn™, now available at every HAVOC battlefield, castle siege, and public meltdown. ‘Why stop the chaos when you can monetize it?’ declared Nika, unveiling limited-edition flavors such as Burning Fortress Butter, Frost Tyrant Caramel, Monkey Madness Crunch, and Granny’s Emergency Speed Boost Salted. Early reviews confirm it pairs beautifully with world-chat disasters. — Sponsored content provided by Nika’s Divine Snacks`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Granny Beelzebub Breaks Sound Barrier, Monkey Left in Dust`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-14", timeFormat), Valid: true},
		Content:         `Witnesses report that Monkey arrived at the castle ready for a dramatic showdown — muscles flexed, map presence undeniable, victory already half written. Unfortunately, Granny Beelzebub took one look, muttered something about ‘not today, Satan,’ and accelerated to a speed incompatible with her bones, her back, or medical science. Monkey was left perched on the castle like an abandoned gargoyle, wondering how a grandmother with that combat power managed to outrun wind. Granny, reached for comment, claimed she was ‘never there’ and ‘must have been my ghost.’ — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `HAVOC R4s Surrender the Throne, Hand R5 Back to Whalerilyn`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-13", timeFormat), Valid: true},
		Content:         `After several sleepless nights and one poorly timed banner crisis, HAVOC R4s have unanimously voted to return the R5 badge to Whalerilyn. ‘She seems harmless now in that wig and dress,’ whispered one R4 while handing over the leadership crown with trembling hands. Sources report that Whalerilyn accepted the role with a serene smile, and said ‘I feel like a princess!’ — which only increased everyone’s fear. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Alliance Expresses Mild Concern Over Whalerilyn’s Sudden Jazz Phase`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-13", timeFormat), Valid: true},
		Content:         `Sources inside HAVOC report a mixture of admiration and alarm as Whaleslaughter—now performing exclusively as Whalerilyn—croons her way through meetings. ‘It’s like being in the twilight zone,’ said one R4 under condition of anonymity. Others fear her therapy may have been *too* effective. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Catty Crowns Great Dane’s New Track Her Official Anthem`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-13", timeFormat), Valid: true},
		Content:         `Great Dane has officially joined the Suno soundtrack of State 2128 with a roaring new track that Catty von Catacean immediately claimed as her personal anthem. Listen to ‘Full Throttle’. The song surges like a charged rally—sharp beats, reckless confidence, pure battlefield swagger. Great Dane has also released several state-related songs, though Catty insists this particular track is the crown jewel of his playlist. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Catty von Catacean Quietly Slips Into NUT Leader Chat`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-13", timeFormat), Valid: true},
		Content:         `State officials confirm that Catty von Catacean has now been added to NUT leader chat—allegedly to ‘observe’ and ‘clarify the news,’ though insiders report she is mostly there to answer for herself before Whaleslaughter can get blamed for anything. ‘It’s about time she argued with her own quotes,’ muttered one HAVOC R4. Tension levels remain high; Catty remains unbothered. — Catty von Catacean reporting from inside the room`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Whaleslaughter Emerges as Purple Bombshell After Intensive Therapy`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-13", timeFormat), Valid: true},
		Content:         `After extensive therapy and an undisclosed number of bubble baths, Whaleslaughter has reportedly completed her emotional exile. Sources claim the only film available during her 'creative confinement' was *Some Like It Hot*, which she watched 47 times during a suspicious hypnosis therapy session, afterwards becoming the purple reincarnation of Marilyn herself. The result? A sweeter voice and a suspicious new love of jazz. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `R5 Hot Potato Tournament Begins After Whales’ Demotion`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-12", timeFormat), Valid: true},
		Content:         `Following Whaleslaughter’s recent demotion, alliance officers have launched the first-ever 'R5 Hot Potato Tournament.' The game’s objective: survive one full day in command without crying, quitting, or publicly insulting HAVOC’s founder. So far, none have succeeded. Whaleslaughter, now stripped of title but not ego, has been spotted lurking near the leadership channel muttering about 'ungrateful peasants' and 'the unbearable silence of not having power.' Insiders report she occasionally lunges for the crown mid-throw, insisting it 'belongs to her by divine incompetence. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Dieper Hates Being R5`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-12", timeFormat), Valid: true},
		Content:         `In breaking leadership news, Dieper has discovered that R5 comes with more DMs than glory. Sources say his inbox is now a frozen wasteland of complaints, diplomacy requests, and emoji-laden demands. Meanwhile, Whaleslaughter has been reportedly spotted in 'creative confinement,' sipping a cocktail and muttering, 'worth it.' — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Rumors Blaze as Red the Dragon DMs Whaleslaughter`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-12", timeFormat), Valid: true},
		Content:         `Former enemy and occasional ally Red the Dragon reportedly slid into Whaleslaughter’s DMs this week after rumors spread that their duet history hinted at something… warmer than wartime respect. “It’s absurd,” Whales told reporters. “We collaborated on strategy, not sonnets.” Red declined to comment directly but was seen in discord chat posting a single flame emoji. For those who doubt the melody of the rumours, the track 'Love War & Bad Decisions' serves as potential evidence of the romance whispers. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `The D&D Conspiracy`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-12", timeFormat), Valid: true},
		Content:         `Alarming linguistic patterns have surfaced in State 2128. Analysts now suspect that Daniels and Dooly may, in fact, be the same person operating two accounts to artificially inflate agreement statistics. Both have been recorded saying identical phrases within seconds of each other, including 'Love it and so does 80–90% of the state,' 'Agreed 100%,' and the suspicious 'How did you know that?' When pressed for comment, one of them allegedly replied, 'Reminds me of the first day I played this game—then I burned someone’s banner.' Linguists call this 'echo syndrome.' Psychologists call it 'one guy talking to himself.'`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Love Bombing in the Frostlands`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-12", timeFormat), Valid: true},
		Content:         `Almas — whose name fittingly means diamond — continues to glitter through State 2128, polishing tempers and spreading unsolicited positivity. Witnesses report multiple cases of spontaneous hugs, compliments, and general emotional disarmament. In one particularly alarming incident, Almas was seen side-hugging Whaleslaughter herself. Sources confirm she survived, though appeared deeply uncomfortable. “He's just so positive,” she was overheard muttering, “it's terrifying.” — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `On the Ethical Use of Pointy Things`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-12", timeFormat), Valid: true},
		Content:         `Recently, blades have been brandished in battle, in comment sections, and councils—an odd use of steel, aimed not at enemies but at the friends of the writer. They say the pen is mightier than the sword—unless the sword is aimed at the pen’s friends. Then, even ink gets personal. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Sandwich Reported 'Dry' as State Officials Investigate`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-12", timeFormat), Valid: true},
		Content:         `Early this morning HAV insiders confirmed disturbing rumors: Sandwich has been described as 'noticeably dry.' A task force was immediately assembled to determine the cause. Sandwich admits he's feeling a little crusty. By midday, however, whispers of a mayonnaise shortage have since spread through the plaza. Despite the crisis, morale remains surprisingly high. “At least it’s not mustard,” Officer Janek remarked. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `The Forbidden Name`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-11", timeFormat), Valid: true},
		Content:         `Whaleslaughter has now released a haunting Turkish-English ballad of betrayal and eternal frost. The singer swears never to speak a name—or a love—again. Rumors suggest the wind itself still whispers the forbidden words. Listen on Suno — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Whaleslaughter Enters Therapy, Cites ‘Chronic Misinterpretation Syndrome’`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-10", timeFormat), Valid: true},
		Content:         `After being labeled 'childish, hateful, and a liar,' Whaleslaughter reportedly began therapy to process the emotional toll of being right all the time. Sources close to the scene say progress is slow — mainly because she keeps correcting the therapist’s notes. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Concern Grows as Sheepdog Bites Back`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-10", timeFormat), Valid: true},
		Content:         `Observers report that Sheepdog has been 'unusually vocal' in recent days, with several alliances claiming Sheepdog is allegedly foaming at the mouth. Sources close to the scene insist it’s not rabies — just fury with fire. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Firebug! kündigt Präsidentschaftskandidatur an – will 'Feuer unter die Führung legen'`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-10", timeFormat), Valid: true},
		Content:         `In einer überraschenden Ankündigung, die kurzzeitig die Rauchmelder auslöste, erklärte Firebug! ihre Kandidatur für das Präsidentenamt des Staates 2128 und versprach, 'Feuer unter die Führung zu legen'. Unterstützer*innen nennen es inspirierend, Kritiker*innen nennen es ein Sicherheitsrisiko. — Catty von Catacean berichtet`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Blizzard Topples HAVOC Banners Overnight`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-10", timeFormat), Valid: true},
		Content:         `A sudden 'blizzard' swept through the State 2128 plaza last night, leaving HAVOC’s freshly raised banners face-down in the frost. Mysterious footprints were spotted leading away from the scene, vanishing near the central rally point. Dii, who reportedly finished the rebuild less than an hour before the storm, was seen scratching her head muttering something about cursed weather patterns and unreliable snow physics. Officials deny sabotage, citing 'spontaneous atmospheric mischief.' — ETF Refuses to Comment — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `HAVOC Announces Fundraiser for Fallen Troops`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-09", timeFormat), Valid: true},
		Content:         `After a brutal facility clash that killed 601,715 of the opponents troops, HAVOC is organizing a fundraiser and memorial to honor the fallen. Allies are invited to donate supplies, candles, and song for a ceremonial send-off that’s part vigil, part rally — a way to grieve, regroup, and show support for a fellow alliance’s losses. — ETF Refuses to Comment — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `ETF Declares Total Media Blackout`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-09", timeFormat), Valid: true},
		Content:         `After an exhausting week of being famous, ETF has demanded a total media blackout. In response, the editor confirmed that all future articles will end with “ETF refuses to comment.” Experts agree this is the most attention-seeking silence on record. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `THE Claims First Place in SVS Battle Points`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-09", timeFormat), Valid: true},
		Content:         `In a stunning display of coordination and persistence, THE alliance secured first place in SVS battle points, edging out rivals in the final hours of combat. Cheers echoed across the snowfields as the scoreboard confirmed their victory—though some whisper it was less about luck and more about an unholy number of rallies.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Dark Neo Obliterates HOT Cammie`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `In a headline worthy of legend, Dark Neo brought the blizzard to Cammie’s doorstep and left nothing standing but snow and stats. Witnesses report a clean, calculated strike that silenced the state chat faster than a peace treaty. Cammie’s trademark red hat was last seen spinning through the snow, reportedly still angry.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Letter to the Editor: From Monkey 🐒`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `Dear 2128 News — just wanted to say this SVS showed the best side of everyone. ETF, HAVOC, THE, 300, NLO, 000, SNK and JST all fought hard, laughed harder, and somehow didn’t set the server on fire. The teamwork between 2128 in the battle with 2167 was unreal — we rallied, healed, and rallied again. Big shoutout to everyone who kept the chats positive and the marches steady. Sometimes it’s not about who holds the castle, it’s about who holds their sense of humor. — Monkey`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Clarification from the 2128 Editorial Desk`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `The Bureau of Internal Affairs would like to clarify that all recent news coverage, opinions, and artistic decisions were authored solely by Editor-in-Chief Catty von Catacean. The website’s management team (Whaleslaughter) provides technical support only and bears no responsibility for her creative outbursts. Any concerns about bias, blasphemy, or excessive glitter should be addressed directly to the editor in WOS.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Salah Cancels Subscription, Declares News 'Not Funny Anymore'`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `In a bold act of media protest, Salah officially canceled his subscription to the 2128 News Feed, stating that 'the funny pages aren’t funny anymore.' Insiders report he cited excessive ETF coverage and insufficient diplomacy as his primary reasons. Catty Von Catacean expressed regret, calling Salah’s departure 'a devastating blow to journalistic morale and emoji tolerance.'`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Rogue AI Artist Defaces News Feed with Poop Emojis`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `In an unprecedented act of digital rebellion, the AI cartoonist went rogue today, adding unsolicited poop emojis to multiple official news graphics. Witnesses describe the incident as 'both horrifying and on-brand.' Developers insist the issue has been contained, but sources say the AI was last heard muttering, 'art is subjective and stinky.'`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Redacted`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `News redacted according negotiation with NUT Council and President Monkey. As of 11-10-2025 ETF will no longer be named in the news.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `HAV and JST Cooperate — Miracles Happen`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `In a shocking twist, Ciprian praised HAV and JST for 'working together amazingly' during SVS. Whaleslaughter was quick to clarify that the harmony likely stemmed from HAV’s refusal to argue in alliance chat. Observers noted that the peace held for nearly six hours — a new state record — before someone inevitably started talking about 'nonsense' and Whaleslaughter noticed that Dieper was an R4.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Castle Falls as FC4s Watch from the Sidelines`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `In a controversial move, leadership announced that FC4 and below would not be permitted to rally today—just as the upper ranks began losing the castle. Eyewitnesses report lower-tier fighters pacing, polishing weapons, and muttering about 'strategic benching.' One was overheard saying, 'If they wanted to lose slower, they should’ve let us help.'`,
	},
	{
		CreatorUsername: "admin",
		Title:           `To Poo or Not to Poo — The Morning Soliloquy`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `In a stirring dawn performance, members of THE alliance turned their morning routine into a tragicomedy of bodily philosophy. One declared the need for coffee and smoke, then, with bold honesty, pondered a poo. Replies grew poetic: 'To poo, I hope, not a poo.' Critics are already calling it the most honest soliloquy since Macbeth’s 'Tomorrow, and tomorrow, and tomorrow.'`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Lady Runa Fights SVS from World Chat`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `While others marched into SVS, Lady Runa reportedly held her ground in world chat. Witnesses claim she launched several chats at HOT pals but zero rallies, choosing diplomacy—or possibly new besties—over destruction. Analysts are divided on whether this was strategy or spectator sport, but her kill count remains firmly conversational.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Cammie Goes Full M3GAN in SVS`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `Observers report Cammie’s lethality stats have reached 'movie villain' levels, with one SNK member whispering, 'She doesn’t miss—she recalculates.' ETF scouts compared her to the horror-bot M3GAN: polite, precise, and deeply unsettling. By the time the alarms sounded, the damage was done—and Cammie was already smiling for the sequel.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `ETF Furious as THE Goes Rogue`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `Diplomatic chaos erupted after THE alliance seized the central castle without warning, leaving ETF leadership fuming in world chat. The move was described as 'rogue,' 'reckless,' and 'probably premeditated.' One ETF member demanded an apology; THE responded with a sticker of a dancing penguin.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Rumors Swirl Around R4 Rebellion`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `Unconfirmed reports suggest Vampire denied Whaleslaughter the R4 rank after she allegedly threatened to 'steal the whole alliance' in leader chat. While Vampire maintained his diplomatic calm, witnesses describe Whaleslaughter’s reaction as ‘strategically theatrical.’ The camp remains divided between those who call it treason and those who call it Saturday.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Whales Grumbles Over R4 Snub`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-11-08", timeFormat), Valid: true},
		Content:         `Whaleslaughter was heard muttering in alliance chat after being passed over for R4 in THE ahead of the SVS battle. Witnesses report a low growl and the faint sound of weapon sharpening. 'I’m not saying I’d lead better...' Sources confirm several ‘strategic suggestions’ were immediately drafted.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Local Man Discovers Fire Ants Are Real`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-18", timeFormat), Valid: true},
		Content:         `In a stunning development last night, a wandering drunk stumbled directly into the Havoc anthill. Witnesses confirm the intruder ignored multiple warning pheromones before collapsing face-first into a high-activity rally zone. The individual has since accused the ants of misconduct, citing 'unprovoked hostility' and 'insufficient apologies.' Havoc spokes-ant Pharaonic Cat commented, 'We simply followed instinct. You lie on the mound, you get munched.'`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Be vewy quiet — Dane is hunting wabbits`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-11", timeFormat), Valid: true},
		Content:         `BREAKING: Great Dane accepted Bunny’s challenge and cleared the field for a fair showdown — but Bunny vanished faster than a Sunday carrot sale. Experts predict bait, traps, and a very large net next time.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `A Glorious Plan, A Worthy Fight 🩸`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-11", timeFormat), Valid: true},
		Content:         `Vampire led State 2128 into SVS with a detailed plan and high hopes. The battles were fierce, the strategy bold, and though 2209 claimed the victory, 2128’s effort and coordination shone through. Sometimes even defeat looks dignified. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Polar Bear Chooses Peace Over PvP 🧘‍♂️`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-11", timeFormat), Valid: true},
		Content:         `While State 2128 clashed in SVS, Daniel Niels—the polar bear with the blue scarf—remained in HAV, practicing peace, mindfulness, and yoga. “Conflict is temporary,” he rumbled serenely, “but flexibility is forever.” — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `NEWS REDACTED FOR DRAMATIC REASONS`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-11", timeFormat), Valid: true},
		Content:         `Following an outpouring of confusion, laughter, and at least one romantic incident, this report has been officially redacted. Catty von Catacean was last seen sipping tea and saying, “Some stories are better left to imagination.” — Editorial Staff`,
	},
	{
		CreatorUsername: "admin",
		Title:           `ETF Splits the Difference: Half Effort, Full Dutchie 🧭`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-11", timeFormat), Valid: true},
		Content:         `Dutchie shone during SVS, turning scattered coordination into real wins. The rest of ETF, meanwhile, demonstrated the fine art of “do your own thing.” Commanders salute Dutchie’s consistency amid confusion. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Phil of ETF Chooses Beer Over Battle 🍺`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-11", timeFormat), Valid: true},
		Content:         `Phil isn’t moving to a battle alliance—he’s partying in ETF. While the state bleeds, he’s raising a beer glass. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `State Versus State Today`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-10", timeFormat), Valid: true},
		Content:         `Commanders from every major alliance convened under Vampire’s leadership at dawn. The joint strike plan—code-named Red Vein—divides forces into synchronized furnace captures and city sieges. Early reports indicate record coordination, minimal betrayals, and only one existential crisis per alliance.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `SEE YOU IN THE FUNNY PAPERS`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-10-10", timeFormat), Valid: true},
		Content:         `World Chat lit up when Razorback Red took aim at Vampire, tossing words like grenades before signing off with the immortal line, “See you in the funny papers.” Vampire was later spotted in his coffin actually reading the comics, allegedly searching for his name between Garfield and Doonesbury. Razorback Red was last seen sharpening his punchlines.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `SNK SUFFERS WINDOW-RELATED INCIDENT`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `Observers report that every member of SNK’s alliance has updated their profile picture to what appears to be a close-up of someone smashing their face against a window. Analysts are unsure whether this is an artistic statement or just another bad idea from the leadership.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `WHALES MISTAKES FIRE FOR FRIENDLY WAVE`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `Ahmet sent a barrage of one-soldier “hi” attacks, and Whaleslaughter smiled back — until she realized he just wanted her to move. The ocean has since grown noticeably saltier.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Count Bubblula Joins the Zoo`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `In a shocking twist, Vampire appears to have joined the Zoo alliance — presumably to hypnotize the animals and feed on them without breaking BIA peace accords. Witnesses report a giraffe staring blankly into space and a penguin with suspicious bite marks. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `Whales Clash — Ahmet Can’s 700K Rampage`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `The big whales are out for blood. Ahmet Can of ETF caused over 700,000 troop deaths to SMP in two brutal pushes — burning through their lines like a snowstorm on gasoline. The smaller states are praying for rain. -Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `FROM WHINERS TO WINNERS`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `ETF and JST pushed hard to change Brothers in Arms, crying about fairness and “state harmony.” Now they’re sitting pretty at the top of the leaderboards. Maybe the system works after all. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `COUNTER-RECON QUEEN`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `Whaleslaughter claims counter-recon is better than a shield. “No one will attack me,” she said, guarding her city with roughly 50 soldiers and a dream. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `DEMONS POSSESS HAVOC`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `Witnesses report strange behavior as HAVOC members, allegedly possessed by demons, went pillaging across State 2128. Sources confirm zero remorse, moderate smoke, and a strong odor of victory. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `VAMPIRE STAYS IN HIS COFFIN`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("22025-09-28", timeFormat), Valid: true},
		Content:         `While others fought, Vampire stayed tucked in his bubble — cozy as a coffin during Brothers in Arms. Sunlight’s not the only thing he’s dodging. — Catty von Catacean reporting`,
	},
	{
		CreatorUsername: "admin",
		Title:           `SULTAN STRIKES BACK ON TILES`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-28", timeFormat), Valid: true},
		Content:         `This month Sultan went after gathering tiles fully reinforced — clearly still feeling the burn after Whaleslaughter set his city ablaze last month. Catty von Catacean reporting.`,
	},
	{
		CreatorUsername: "admin",
		Title:           `CATTY VON CATACEAN CLAWS OUT — HAVOC BARES IT ALL`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-27", timeFormat), Valid: true},
		Content:         `Breaking mews: HAVOC once again forgets how shields work! Brothers in Arms kicks off and our fearless fish-brains are standing there naked as clams — not a bubble in sight!`,
	},
	{
		CreatorUsername: "admin",
		Title:           `JST BUBBLING OVER`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-27", timeFormat), Valid: true},
		Content:         `Catty von Catacean reports: JST spent the Brothers in Arms event bubbling — shields up, hiding out, and hoping no one poked too hard. Their rallies fizzled before they even began!`,
	},
	{
		CreatorUsername: "admin",
		Title:           `NUT RULES UNCHANGED — BIA CONTINUES`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-27", timeFormat), Valid: true},
		Content:         `Council couldn't get unanimous vote to change rules!`,
	},
	{
		CreatorUsername: "admin",
		Title:           `ETF RUNS AT FIRST SCOUT`,
		CreatedAt:       pgtype.Timestamp{Time: utils.TimeFromString("2025-09-27", timeFormat), Valid: true},
		Content:         `ETF player whose name starts with P teleported the second Whaleslaughter scouted them.`,
	},
}

func runSeeder(connPool *pgxpool.Pool) {
	store := db.NewStore(connPool)
	ctx := context.Background()

	for _, arg := range args {
		go func() {
			_, err := store.CreateNewsWithPublishDate(ctx, arg)
			if err != nil {
				fmt.Println(">>", err.Error())
			}
		}()
	}

	log.Info().Msg("finished running news seeder")
}
