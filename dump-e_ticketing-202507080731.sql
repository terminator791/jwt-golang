--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Debian 16.9-1.pgdg120+1)
-- Dumped by pg_dump version 17.5

-- Started on 2025-07-08 07:31:47

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- TOC entry 3466 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 223 (class 1259 OID 16988)
-- Name: card_balance_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.card_balance_logs (
    log_id uuid DEFAULT gen_random_uuid() NOT NULL,
    card_id uuid NOT NULL,
    previous_balance numeric(12,2) NOT NULL,
    current_balance numeric(12,2) NOT NULL,
    amount_changed numeric(12,2) NOT NULL,
    change_type character varying(20) NOT NULL,
    transaction_id uuid,
    logged_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    description character varying(255)
);


ALTER TABLE public.card_balance_logs OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16945)
-- Name: cards; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cards (
    card_id uuid DEFAULT gen_random_uuid() NOT NULL,
    card_number character varying(20) NOT NULL,
    user_id uuid,
    current_balance numeric(12,2) DEFAULT 0 NOT NULL,
    status character varying(10) DEFAULT 'ACTIVE'::character varying NOT NULL,
    issued_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    is_active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.cards OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 16935)
-- Name: fare_matrices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.fare_matrices (
    fare_id uuid DEFAULT gen_random_uuid() NOT NULL,
    from_terminal_id bigint NOT NULL,
    to_terminal_id bigint NOT NULL,
    base_fare numeric(10,2) NOT NULL,
    peak_hour_multiplier numeric(4,2) DEFAULT 1,
    effective_from timestamp with time zone NOT NULL,
    effective_until timestamp with time zone,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.fare_matrices OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 16957)
-- Name: gates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.gates (
    gate_id uuid DEFAULT gen_random_uuid() NOT NULL,
    terminal_id bigint NOT NULL,
    gate_code character varying(20) NOT NULL,
    gate_type character varying(10) NOT NULL,
    status character varying(20) DEFAULT 'ACTIVE'::character varying NOT NULL,
    ip_address character varying(45),
    last_heartbeat timestamp with time zone,
    is_online boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.gates OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 16997)
-- Name: sync_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sync_logs (
    sync_id uuid DEFAULT gen_random_uuid() NOT NULL,
    terminal_id bigint NOT NULL,
    sync_type character varying(20) NOT NULL,
    sync_data json,
    sync_status character varying(15) NOT NULL,
    sync_started_at timestamp with time zone NOT NULL,
    sync_completed_at timestamp with time zone,
    error_message character varying(255),
    retry_count bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.sync_logs OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16927)
-- Name: terminals; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.terminals (
    terminal_id bigint NOT NULL,
    terminal_name character varying(100) NOT NULL,
    terminal_code character varying(20) NOT NULL,
    location character varying(255),
    latitude numeric(10,7),
    longitude numeric(10,7),
    is_active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.terminals OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 16926)
-- Name: terminals_terminal_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.terminals_terminal_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.terminals_terminal_id_seq OWNER TO postgres;

--
-- TOC entry 3467 (class 0 OID 0)
-- Dependencies: 216
-- Name: terminals_terminal_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.terminals_terminal_id_seq OWNED BY public.terminals.terminal_id;


--
-- TOC entry 222 (class 1259 OID 16981)
-- Name: top_ups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.top_ups (
    top_up_id uuid DEFAULT gen_random_uuid() NOT NULL,
    card_id uuid NOT NULL,
    amount numeric(12,2) NOT NULL,
    payment_method character varying(20) NOT NULL,
    payment_reference character varying(100),
    status character varying(10) NOT NULL,
    processed_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.top_ups OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 16967)
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    transaction_id uuid DEFAULT gen_random_uuid() NOT NULL,
    card_id uuid NOT NULL,
    origin_terminal_id bigint,
    destination_terminal_id bigint,
    checkin_gate_id uuid,
    checkout_gate_id uuid,
    checkin_time timestamp with time zone NOT NULL,
    checkout_time timestamp with time zone,
    fare_amount numeric(10,2),
    balance_before numeric(12,2) NOT NULL,
    balance_after numeric(12,2),
    transaction_status character varying(20) NOT NULL,
    transaction_type character varying(20) DEFAULT 'REGULAR'::character varying NOT NULL,
    reference_number character varying(50) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    is_synced boolean DEFAULT false NOT NULL
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 16915)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    user_id uuid DEFAULT gen_random_uuid() NOT NULL,
    full_name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password character varying(255) NOT NULL,
    phone character varying(20),
    date_of_birth timestamp with time zone,
    user_type character varying(20) DEFAULT 'CUSTOMER'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 3237 (class 2604 OID 16930)
-- Name: terminals terminal_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.terminals ALTER COLUMN terminal_id SET DEFAULT nextval('public.terminals_terminal_id_seq'::regclass);


--
-- TOC entry 3459 (class 0 OID 16988)
-- Dependencies: 223
-- Data for Name: card_balance_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.card_balance_logs (log_id, card_id, previous_balance, current_balance, amount_changed, change_type, transaction_id, logged_at, description) FROM stdin;
\.


--
-- TOC entry 3455 (class 0 OID 16945)
-- Dependencies: 219
-- Data for Name: cards; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cards (card_id, card_number, user_id, current_balance, status, issued_at, expires_at, created_at, updated_at, is_active) FROM stdin;
\.


--
-- TOC entry 3454 (class 0 OID 16935)
-- Dependencies: 218
-- Data for Name: fare_matrices; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.fare_matrices (fare_id, from_terminal_id, to_terminal_id, base_fare, peak_hour_multiplier, effective_from, effective_until, is_active, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3456 (class 0 OID 16957)
-- Dependencies: 220
-- Data for Name: gates; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.gates (gate_id, terminal_id, gate_code, gate_type, status, ip_address, last_heartbeat, is_online, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3460 (class 0 OID 16997)
-- Dependencies: 224
-- Data for Name: sync_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sync_logs (sync_id, terminal_id, sync_type, sync_data, sync_status, sync_started_at, sync_completed_at, error_message, retry_count) FROM stdin;
\.


--
-- TOC entry 3453 (class 0 OID 16927)
-- Dependencies: 217
-- Data for Name: terminals; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.terminals (terminal_id, terminal_name, terminal_code, location, latitude, longitude, is_active, created_at, updated_at) FROM stdin;
1	Terminal Pusat	TPT01	Jakarta Pusat	-6.1751100	106.8650360	t	2025-07-07 17:47:03.66744+00	2025-07-07 17:47:03.66744+00
2	Terminal Pusat	TPT02	Jakarta Pusat	-6.1751100	106.8650360	t	2025-07-07 17:50:13.145905+00	2025-07-07 17:50:13.145905+00
\.


--
-- TOC entry 3458 (class 0 OID 16981)
-- Dependencies: 222
-- Data for Name: top_ups; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.top_ups (top_up_id, card_id, amount, payment_method, payment_reference, status, processed_at, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3457 (class 0 OID 16967)
-- Dependencies: 221
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (transaction_id, card_id, origin_terminal_id, destination_terminal_id, checkin_gate_id, checkout_gate_id, checkin_time, checkout_time, fare_amount, balance_before, balance_after, transaction_status, transaction_type, reference_number, created_at, updated_at, is_synced) FROM stdin;
\.


--
-- TOC entry 3451 (class 0 OID 16915)
-- Dependencies: 215
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (user_id, full_name, email, password, phone, date_of_birth, user_type, created_at, updated_at) FROM stdin;
2a5db352-06d9-4a64-9f3c-4f5b075206db	Admin Sistem E-Ticketing	admin@e-ticketing.com	$2a$10$JyYPFKkjq3mHbbPCg6WpmemIJ6wKw5zUOYLctKd450FIhrrdNwhji	081234567890	1990-01-01 00:00:00+00	ADMIN	2025-07-07 17:42:06.552367+00	2025-07-07 17:42:06.552367+00
a08bd481-60b7-4719-9641-6c43aea2380e	Budi Santoso	budi@example.com	$2a$10$XWZJOV4z9VWgvR7qemu1n.Fnj7XptwbfQZdP5TwTPU5xLBnn5W0mi	0812345678901	1990-01-15 00:00:00+00	CUSTOMER	2025-07-07 17:49:17.367898+00	2025-07-07 17:49:17.367898+00
\.


--
-- TOC entry 3468 (class 0 OID 0)
-- Dependencies: 216
-- Name: terminals_terminal_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.terminals_terminal_id_seq', 2, true);


--
-- TOC entry 3289 (class 2606 OID 16994)
-- Name: card_balance_logs card_balance_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.card_balance_logs
    ADD CONSTRAINT card_balance_logs_pkey PRIMARY KEY (log_id);


--
-- TOC entry 3270 (class 2606 OID 16954)
-- Name: cards cards_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_pkey PRIMARY KEY (card_id);


--
-- TOC entry 3266 (class 2606 OID 16942)
-- Name: fare_matrices fare_matrices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fare_matrices
    ADD CONSTRAINT fare_matrices_pkey PRIMARY KEY (fare_id);


--
-- TOC entry 3274 (class 2606 OID 16964)
-- Name: gates gates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gates
    ADD CONSTRAINT gates_pkey PRIMARY KEY (gate_id);


--
-- TOC entry 3294 (class 2606 OID 17005)
-- Name: sync_logs sync_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sync_logs
    ADD CONSTRAINT sync_logs_pkey PRIMARY KEY (sync_id);


--
-- TOC entry 3264 (class 2606 OID 16933)
-- Name: terminals terminals_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.terminals
    ADD CONSTRAINT terminals_pkey PRIMARY KEY (terminal_id);


--
-- TOC entry 3287 (class 2606 OID 16986)
-- Name: top_ups top_ups_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.top_ups
    ADD CONSTRAINT top_ups_pkey PRIMARY KEY (top_up_id);


--
-- TOC entry 3284 (class 2606 OID 16974)
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id);


--
-- TOC entry 3261 (class 2606 OID 16923)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- TOC entry 3290 (class 1259 OID 16996)
-- Name: idx_card_balance_logs_card_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_card_balance_logs_card_id ON public.card_balance_logs USING btree (card_id);


--
-- TOC entry 3291 (class 1259 OID 16995)
-- Name: idx_card_balance_logs_transaction_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_card_balance_logs_transaction_id ON public.card_balance_logs USING btree (transaction_id);


--
-- TOC entry 3271 (class 1259 OID 16956)
-- Name: idx_cards_card_number; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_cards_card_number ON public.cards USING btree (card_number);


--
-- TOC entry 3272 (class 1259 OID 16955)
-- Name: idx_cards_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_cards_user_id ON public.cards USING btree (user_id);


--
-- TOC entry 3267 (class 1259 OID 16944)
-- Name: idx_fare_matrices_from_terminal_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_fare_matrices_from_terminal_id ON public.fare_matrices USING btree (from_terminal_id);


--
-- TOC entry 3268 (class 1259 OID 16943)
-- Name: idx_fare_matrices_to_terminal_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_fare_matrices_to_terminal_id ON public.fare_matrices USING btree (to_terminal_id);


--
-- TOC entry 3275 (class 1259 OID 16965)
-- Name: idx_gates_gate_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_gates_gate_code ON public.gates USING btree (gate_code);


--
-- TOC entry 3276 (class 1259 OID 16966)
-- Name: idx_gates_terminal_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_gates_terminal_id ON public.gates USING btree (terminal_id);


--
-- TOC entry 3292 (class 1259 OID 17006)
-- Name: idx_sync_logs_terminal_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sync_logs_terminal_id ON public.sync_logs USING btree (terminal_id);


--
-- TOC entry 3262 (class 1259 OID 16934)
-- Name: idx_terminals_terminal_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_terminals_terminal_code ON public.terminals USING btree (terminal_code);


--
-- TOC entry 3285 (class 1259 OID 16987)
-- Name: idx_top_ups_card_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_top_ups_card_id ON public.top_ups USING btree (card_id);


--
-- TOC entry 3277 (class 1259 OID 16980)
-- Name: idx_transactions_card_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_card_id ON public.transactions USING btree (card_id);


--
-- TOC entry 3278 (class 1259 OID 16977)
-- Name: idx_transactions_checkin_gate_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_checkin_gate_id ON public.transactions USING btree (checkin_gate_id);


--
-- TOC entry 3279 (class 1259 OID 16976)
-- Name: idx_transactions_checkout_gate_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_checkout_gate_id ON public.transactions USING btree (checkout_gate_id);


--
-- TOC entry 3280 (class 1259 OID 16978)
-- Name: idx_transactions_destination_terminal_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_destination_terminal_id ON public.transactions USING btree (destination_terminal_id);


--
-- TOC entry 3281 (class 1259 OID 16979)
-- Name: idx_transactions_origin_terminal_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_transactions_origin_terminal_id ON public.transactions USING btree (origin_terminal_id);


--
-- TOC entry 3282 (class 1259 OID 16975)
-- Name: idx_transactions_reference_number; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_transactions_reference_number ON public.transactions USING btree (reference_number);


--
-- TOC entry 3258 (class 1259 OID 16925)
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- TOC entry 3259 (class 1259 OID 16924)
-- Name: idx_users_phone; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_users_phone ON public.users USING btree (phone);


--
-- TOC entry 3305 (class 2606 OID 17047)
-- Name: card_balance_logs fk_cards_card_balance_logs; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.card_balance_logs
    ADD CONSTRAINT fk_cards_card_balance_logs FOREIGN KEY (card_id) REFERENCES public.cards(card_id) ON DELETE CASCADE;


--
-- TOC entry 3304 (class 2606 OID 17042)
-- Name: top_ups fk_cards_top_ups; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.top_ups
    ADD CONSTRAINT fk_cards_top_ups FOREIGN KEY (card_id) REFERENCES public.cards(card_id) ON DELETE CASCADE;


--
-- TOC entry 3299 (class 2606 OID 17017)
-- Name: transactions fk_cards_transactions; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_cards_transactions FOREIGN KEY (card_id) REFERENCES public.cards(card_id) ON DELETE CASCADE;


--
-- TOC entry 3300 (class 2606 OID 17032)
-- Name: transactions fk_gates_checkin_transactions; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_gates_checkin_transactions FOREIGN KEY (checkin_gate_id) REFERENCES public.gates(gate_id) ON DELETE SET NULL;


--
-- TOC entry 3301 (class 2606 OID 17037)
-- Name: transactions fk_gates_checkout_transactions; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_gates_checkout_transactions FOREIGN KEY (checkout_gate_id) REFERENCES public.gates(gate_id) ON DELETE SET NULL;


--
-- TOC entry 3302 (class 2606 OID 17027)
-- Name: transactions fk_terminals_dest_transactions; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_terminals_dest_transactions FOREIGN KEY (destination_terminal_id) REFERENCES public.terminals(terminal_id) ON DELETE SET NULL;


--
-- TOC entry 3295 (class 2606 OID 17062)
-- Name: fare_matrices fk_terminals_from_fare_matrices; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fare_matrices
    ADD CONSTRAINT fk_terminals_from_fare_matrices FOREIGN KEY (from_terminal_id) REFERENCES public.terminals(terminal_id) ON DELETE CASCADE;


--
-- TOC entry 3298 (class 2606 OID 17012)
-- Name: gates fk_terminals_gates; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gates
    ADD CONSTRAINT fk_terminals_gates FOREIGN KEY (terminal_id) REFERENCES public.terminals(terminal_id) ON DELETE CASCADE;


--
-- TOC entry 3303 (class 2606 OID 17022)
-- Name: transactions fk_terminals_origin_transactions; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_terminals_origin_transactions FOREIGN KEY (origin_terminal_id) REFERENCES public.terminals(terminal_id) ON DELETE SET NULL;


--
-- TOC entry 3307 (class 2606 OID 17057)
-- Name: sync_logs fk_terminals_sync_logs; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sync_logs
    ADD CONSTRAINT fk_terminals_sync_logs FOREIGN KEY (terminal_id) REFERENCES public.terminals(terminal_id) ON DELETE CASCADE;


--
-- TOC entry 3296 (class 2606 OID 17067)
-- Name: fare_matrices fk_terminals_to_fare_matrices; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fare_matrices
    ADD CONSTRAINT fk_terminals_to_fare_matrices FOREIGN KEY (to_terminal_id) REFERENCES public.terminals(terminal_id) ON DELETE CASCADE;


--
-- TOC entry 3306 (class 2606 OID 17052)
-- Name: card_balance_logs fk_transactions_card_balance_logs; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.card_balance_logs
    ADD CONSTRAINT fk_transactions_card_balance_logs FOREIGN KEY (transaction_id) REFERENCES public.transactions(transaction_id) ON DELETE SET NULL;


--
-- TOC entry 3297 (class 2606 OID 17007)
-- Name: cards fk_users_cards; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT fk_users_cards FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE SET NULL;


-- Completed on 2025-07-08 07:31:47

--
-- PostgreSQL database dump complete
--

