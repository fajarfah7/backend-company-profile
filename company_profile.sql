-- Adminer 4.7.8 PostgreSQL dump

DROP TABLE IF EXISTS "cart_products";
DROP SEQUENCE IF EXISTS cart_products_id_seq;
CREATE SEQUENCE cart_products_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."cart_products" (
    "id" integer DEFAULT nextval('cart_products_id_seq') NOT NULL,
    "cart_id" bigint NOT NULL,
    "product_id" bigint NOT NULL,
    "amount" integer NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp
) WITH (oids = false);

INSERT INTO "cart_products" ("id", "cart_id", "product_id", "amount", "created_at", "updated_at") VALUES
(26,	29,	12,	2,	'2021-11-29 14:03:45.497057',	'2021-11-29 14:39:59.720852'),
(28,	31,	12,	1,	'2021-11-29 14:51:38.500518',	'2021-11-29 14:51:38.500519'),
(29,	32,	13,	1,	'2021-11-29 14:57:03.7584',	'2021-11-29 14:57:03.7584'),
(30,	33,	9,	2,	'2021-11-29 15:02:37.645897',	'2021-11-29 15:02:39.819116'),
(31,	34,	13,	1,	'2021-11-29 15:03:05.940743',	'2021-11-29 15:03:05.940743'),
(11,	30,	12,	1,	'2021-11-28 09:59:20.408795',	'2021-11-28 09:59:20.408795'),
(12,	30,	13,	1,	'2021-11-28 09:59:32.513734',	'2021-11-28 09:59:32.513734'),
(37,	36,	13,	1,	'2021-12-01 20:06:14.245911',	'2021-12-01 20:06:14.245911'),
(33,	35,	12,	1,	'2021-11-29 15:04:28.324367',	'2021-12-01 20:38:39.129453'),
(38,	36,	12,	3,	'2021-12-01 20:06:18.875745',	'2021-12-01 20:40:38.887649'),
(39,	37,	12,	1,	'2021-12-01 20:41:12.542057',	'2021-12-01 20:41:12.542057'),
(40,	37,	11,	2,	'2021-12-01 20:45:50.395483',	'2021-12-01 20:46:06.453843'),
(36,	35,	11,	1,	'2021-11-29 15:15:09.284895',	'2021-12-03 20:05:09.17891'),
(35,	35,	13,	2,	'2021-11-29 15:14:54.567951',	'2021-12-03 20:05:10.440917'),
(34,	35,	9,	3,	'2021-11-29 15:05:10.36083',	'2021-12-03 20:05:13.163191'),
(41,	38,	12,	1,	'2021-12-03 20:05:57.463648',	'2021-12-03 20:05:57.463648'),
(42,	38,	14,	1,	'2021-12-14 14:50:35.512293',	'2021-12-14 14:50:35.512293'),
(27,	29,	10,	3,	'2021-11-29 14:11:21.905877',	'2021-11-29 14:12:01.246388'),
(43,	38,	9,	1,	'2021-12-14 14:50:53.158113',	'2021-12-14 14:50:53.158113');

DROP TABLE IF EXISTS "carts";
DROP SEQUENCE IF EXISTS user_carts_id_seq;
CREATE SEQUENCE user_carts_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."carts" (
    "id" integer DEFAULT nextval('user_carts_id_seq') NOT NULL,
    "user_id" integer NOT NULL,
    "status" smallint NOT NULL,
    "payment_code" character varying(64),
    "shipment_code" character varying(20),
    "received_at" timestamp
) WITH (oids = false);

COMMENT ON COLUMN "public"."carts"."status" IS '0=not paid yet, 1=paid, 2=sent, 3=received and completed';

INSERT INTO "carts" ("id", "user_id", "status", "payment_code", "shipment_code", "received_at") VALUES
(36,	5,	2,	'19ca14e7ea6328a42e0eb13d585e4c22',	'asdf',	NULL),
(35,	4,	1,	'1c383cd30b7c298ab50293adfecb7b18',	NULL,	NULL),
(38,	4,	0,	NULL,	NULL,	NULL),
(34,	4,	3,	'e369853df766fa44e1ed0ff613f563bd',	'ddd',	'2021-12-03 20:40:17.067495'),
(29,	4,	1,	'6ea9ab1baa0efb9e19094440c317e21b',	NULL,	NULL),
(31,	4,	1,	'c16a5320fa475530d9583c34fd356ef5',	NULL,	NULL),
(32,	4,	1,	'6364d3f0f495b6ab9dcf8d3b5c6e0b01',	NULL,	NULL),
(33,	4,	1,	'182be0c5cdcd5072bb1864cdee4d3d6e',	NULL,	NULL),
(37,	5,	0,	NULL,	NULL,	NULL);

DROP TABLE IF EXISTS "categories";
DROP SEQUENCE IF EXISTS categories_id_seq;
CREATE SEQUENCE categories_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."categories" (
    "id" integer DEFAULT nextval('categories_id_seq') NOT NULL,
    "name" character varying(30) NOT NULL
) WITH (oids = false);

INSERT INTO "categories" ("id", "name") VALUES
(7,	'Gundam Wing Endless Waltz'),
(9,	'Gundam Seed');

DROP TABLE IF EXISTS "images";
DROP SEQUENCE IF EXISTS images_id_seq;
CREATE SEQUENCE images_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."images" (
    "id" integer DEFAULT nextval('images_id_seq') NOT NULL,
    "name" character varying(30) NOT NULL,
    "path" character varying(255) NOT NULL,
    "type" character varying(15) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp
) WITH (oids = false);


DROP TABLE IF EXISTS "products";
DROP SEQUENCE IF EXISTS products_id_seq;
CREATE SEQUENCE products_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."products" (
    "id" integer DEFAULT nextval('products_id_seq') NOT NULL,
    "uom_id" integer NOT NULL,
    "name" character varying(30) NOT NULL,
    "description" text NOT NULL,
    "dimension_description" character varying(255) NOT NULL,
    "stock" integer NOT NULL,
    "price" numeric(10,2) NOT NULL,
    "discount" smallint NOT NULL,
    "final_price" numeric(10,2) NOT NULL,
    "expired_at" date,
    "created_at" timestamp,
    "updated_at" timestamp,
    "product_image" character varying(255) NOT NULL,
    "status" boolean NOT NULL,
    "is_have_expiry" boolean NOT NULL,
    "category_id" integer
) WITH (oids = false);

INSERT INTO "products" ("id", "uom_id", "name", "description", "dimension_description", "stock", "price", "discount", "final_price", "expired_at", "created_at", "updated_at", "product_image", "status", "is_have_expiry", "category_id") VALUES
(14,	1,	'Gundam Strike Freedom',	'Gundam Strike Freedom from Gundam Seed',	'width: 100cm, height 100cm',	2,	100.00,	0,	0.00,	NULL,	'2021-12-06 20:12:23.351936',	'2021-12-13 16:54:27.741692',	'gundam_strike_freedom.jpeg',	't',	'f',	9),
(15,	1,	'Gundam Justice',	'Gundam Justice from Gundam Seed',	'w: 10, h:100',	4,	100.00,	0,	0.00,	NULL,	'2021-12-13 16:55:58.063226',	'2021-12-13 16:56:28.296413',	'gundam_justice.jpeg',	't',	'f',	9),
(9,	1,	'Gundam Wing Zero',	'This is gundam wing zero from gundam wing',	'width 10cm, height 50cm',	10,	100.00,	0,	0.00,	NULL,	'2021-11-24 17:47:46.663162',	'2021-11-24 17:47:46.663162',	'gundam_wing.jpg',	't',	'f',	7),
(11,	1,	'Gundam Altron',	'This is gundam altron from gundam wing',	'width 10cm, height 50cm',	8,	100.00,	0,	0.00,	NULL,	'2021-11-24 17:50:42.664227',	'2021-11-24 17:50:42.664227',	'gundam_altron.jpg',	't',	'f',	7),
(12,	1,	'Gundam Heavyarms',	'This is gundam heavyarms from gundam wing',	'width 10, height 50',	7,	100.00,	0,	0.00,	NULL,	'2021-11-24 17:51:53.031222',	'2021-11-24 17:51:53.031222',	'gundam_heavyarms.jpeg',	't',	'f',	7),
(10,	1,	'Gundam Deathscythe Hell',	'This is gundam deathscythe hell from gundam wing',	'width 10cm, height 50cm',	5,	100.00,	0,	0.00,	NULL,	'2021-11-24 17:49:19.518584',	'2021-11-28 13:37:30.706204',	'gundam_deathscythe_hell.jpg',	't',	'f',	7),
(13,	1,	'Gundam Sandrock',	'This is gundam sandrock from gundam wing',	'width 10, height 50',	8,	100.00,	0,	0.00,	NULL,	'2021-11-24 17:52:42.155201',	'2021-12-06 20:11:18.132666',	'gundam_sandrock.jpeg',	't',	'f',	7);

DROP TABLE IF EXISTS "texts";
DROP SEQUENCE IF EXISTS texts_id_seq;
CREATE SEQUENCE texts_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."texts" (
    "id" integer DEFAULT nextval('texts_id_seq') NOT NULL,
    "image_id" integer,
    "name" character varying(30) NOT NULL,
    "type" character varying(20) NOT NULL,
    "text" text NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp
) WITH (oids = false);

INSERT INTO "texts" ("id", "image_id", "name", "type", "text", "created_at", "updated_at") VALUES
(5,	NULL,	'Goal',	'about_us_goal',	'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using ''Content here, content here'', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for ''lorem ipsum'' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).',	NULL,	NULL),
(6,	NULL,	'Achievement',	'about_us_achievement',	'Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.',	NULL,	NULL),
(4,	NULL,	'History',	'about_us_history',	'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry''s standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.',	NULL,	NULL);

DROP TABLE IF EXISTS "uoms";
DROP SEQUENCE IF EXISTS measurements_id_seq;
CREATE SEQUENCE measurements_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."uoms" (
    "id" integer DEFAULT nextval('measurements_id_seq') NOT NULL,
    "name" character varying(25) NOT NULL,
    "amount" numeric(10,2) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp
) WITH (oids = false);

INSERT INTO "uoms" ("id", "name", "amount", "created_at", "updated_at") VALUES
(1,	'piece',	1.00,	'2021-10-28 16:23:05.367387',	'2021-10-28 16:23:05.367387'),
(2,	'dozen',	12.00,	'2021-10-28 16:23:24.849243',	'2021-10-28 16:23:24.849243'),
(3,	'score',	20.00,	'2021-10-28 16:23:35.912693',	'2021-10-28 16:23:35.912693');

DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_id_seq;
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "name" character varying(30) NOT NULL,
    "email" character varying(30) NOT NULL,
    "username" character varying(30) NOT NULL,
    "password" character varying(255) NOT NULL,
    "token" text,
    "token_expired" date,
    "address" text NOT NULL,
    "phone_number" character varying(15) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp
) WITH (oids = false);

INSERT INTO "users" ("id", "name", "email", "username", "password", "token", "token_expired", "address", "phone_number", "created_at", "updated_at") VALUES
(4,	'Fajar F',	'fajar@mail.com',	'fajar',	'$2a$10$.oO8WohO1h7TD4LZeV.AWeKr3TvIFJUL4705Kbh4Rfn2rRZDClpJm',	NULL,	NULL,	'Desa Blimbing, RT01/01, Kecamatan Mandiraja, Banjarnegara, Jawa Tengah, Indonesia',	'081229439753',	'2021-11-25 20:49:11.657721',	'2021-12-10 10:47:59.313749'),
(8,	'aaaa',	'aaaa@mail.com',	'aaaa',	'$2a$10$2jKAxOgSQQQ.WyZASOZRtePRT/yXZNvbpfAPJ7F78gQotuKzpKCIK',	NULL,	NULL,	'aaaa',	'1234',	'2021-12-11 22:09:13.407349',	'2021-12-11 22:09:13.407349'),
(7,	'test',	'test@mail.com',	'test',	'$2a$10$Y5N0oVg9yEwRF8Tb7C687e8rif.qPcmMitTNiZRDKm/oMfK4LkaN.',	NULL,	NULL,	'test address',	'12345',	'2021-11-28 09:58:55.994823',	'2021-11-28 09:58:55.994823'),
(5,	'Fahrurozi',	'fahrurozi@mail.com',	'fahrurozi',	'$2a$10$IX1sbJCzd3mE2FQOxNTZeuYMwUXJhOF2bBUFyQsXiP0U92.oNvwzi',	NULL,	NULL,	'Some place in this world',	'085700510051',	'2021-11-25 20:53:28.953338',	'2021-11-25 20:53:28.953338');

-- 2021-12-18 19:03:06.172696+07
