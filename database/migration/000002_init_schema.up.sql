INSERT INTO public."Users" VALUES (1, 'karl', 'karlkarl', 'karl@karl.karl', false, '2023-03-11 21:12:11.767219+00', '2023-03-11 21:12:11.767219+00', '', false);
INSERT INTO public."Users" VALUES (2, 'luke', 'lukeluke', 'luke@luke.net', true, '2023-03-12 15:02:13.984207+00', '2023-03-12 15:02:13.984207+00', 'aaaaaabcd', false);

INSERT INTO public."AdminAssignments" VALUES (1, 1, 1, '2023-03-12 15:00:14.597256+00', false);

INSERT INTO public."GlobalModeratorAssignments" VALUES (1, 2, '2023-03-12 18:50:58.547691+00', 1, false);

INSERT INTO public."Domains" VALUES ('abc.com');

INSERT INTO public."Paths" VALUES (1, 'abc.com', '');

INSERT INTO public."Comments" VALUES (1, 1, 2, 'Hey, ABC!', '2023-03-12 16:12:56.953205+00', NULL, false, false);
INSERT INTO public."Comments" VALUES (2, 1, 1, 'Sup, luke.', '2023-03-12 16:13:31.931985+00', 1, false, false);

INSERT INTO public."VoteRecords" VALUES (1, 'factual', 1, 1);
INSERT INTO public."VoteRecords" VALUES (1, 'agree', 1, 1);
INSERT INTO public."VoteRecords" VALUES (1, 'funny', 1, 1);
INSERT INTO public."VoteRecords" VALUES (1, 'factual', 2, 1);
INSERT INTO public."VoteRecords" VALUES (1, 'funny', 2, 1);
