--
-- PostgreSQL database dump
--

\restrict gA1daJdzDra9ovXP0herlMbqmJWuZ7zPO9h9FI1CNztDdQIVUHnKH9zCcOASz3w

-- Dumped from database version 18.0 (Ubuntu 18.0-1.pgdg24.04+3)
-- Dumped by pg_dump version 18.0 (Ubuntu 18.0-1.pgdg24.04+3)

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
-- Name: public; Type: SCHEMA; Schema: -; Owner: daybook_user
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO daybook_user;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: account_types; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.account_types (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    icon text,
    description text,
    active boolean DEFAULT true,
    sort_order bigint DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.account_types OWNER TO daybook_user;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.accounts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    type text NOT NULL,
    initial_balance numeric DEFAULT 0,
    balance numeric DEFAULT 0,
    currency text DEFAULT 'BDT'::text,
    description text,
    institution text,
    account_number text,
    last_reconciled timestamp with time zone,
    reconciliation_difference numeric DEFAULT 0,
    active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.accounts OWNER TO daybook_user;

--
-- Name: activity_logs; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.activity_logs (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    action text NOT NULL,
    module text NOT NULL,
    entity_type text,
    entity_id uuid,
    description text,
    ip_address text,
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.activity_logs OWNER TO daybook_user;

--
-- Name: asset_attachments; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.asset_attachments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    asset_id uuid NOT NULL,
    file_name text NOT NULL,
    original_name text NOT NULL,
    file_path text NOT NULL,
    file_url text NOT NULL,
    file_size bigint,
    mime_type text,
    attachment_type text,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.asset_attachments OWNER TO daybook_user;

--
-- Name: assets; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.assets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    description text,
    category text,
    brand text,
    model text,
    serial_number text,
    purchase_date timestamp without time zone NOT NULL,
    purchase_price numeric NOT NULL,
    purchase_location text,
    warranty_start_date timestamp without time zone,
    warranty_end_date timestamp without time zone,
    warranty_provider text,
    warranty_type text,
    status text NOT NULL,
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.assets OWNER TO daybook_user;

--
-- Name: bill_payments; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.bill_payments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    bill_id uuid NOT NULL,
    amount numeric NOT NULL,
    payment_date timestamp with time zone NOT NULL,
    account_id uuid,
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.bill_payments OWNER TO daybook_user;

--
-- Name: bills; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.bills (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    category text NOT NULL,
    amount numeric NOT NULL,
    frequency text NOT NULL,
    start_date timestamp with time zone NOT NULL,
    due_day bigint,
    last_paid_date timestamp with time zone,
    last_paid_amount numeric DEFAULT 0,
    auto_pay boolean DEFAULT false,
    reminder_days bigint DEFAULT 3,
    active boolean DEFAULT true,
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.bills OWNER TO daybook_user;

--
-- Name: budgets; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.budgets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    category_id text NOT NULL,
    amount numeric NOT NULL,
    period text NOT NULL,
    custom_start_date timestamp with time zone,
    custom_end_date timestamp with time zone,
    rollover boolean DEFAULT false,
    alert_threshold numeric DEFAULT 80,
    enabled boolean DEFAULT true,
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.budgets OWNER TO daybook_user;

--
-- Name: credit_card_payments; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.credit_card_payments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    account_id uuid NOT NULL,
    amount numeric NOT NULL,
    payment_date timestamp with time zone NOT NULL,
    description text,
    transaction_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.credit_card_payments OWNER TO daybook_user;

--
-- Name: credit_card_transactions; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.credit_card_transactions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    transaction_id uuid,
    category_id text,
    amount numeric NOT NULL,
    description text,
    merchant text,
    date timestamp with time zone NOT NULL,
    type text NOT NULL,
    tags jsonb,
    attachments jsonb,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.credit_card_transactions OWNER TO daybook_user;

--
-- Name: credit_cards; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.credit_cards (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    last_four_digits text,
    card_network text,
    credit_limit numeric NOT NULL,
    current_balance numeric DEFAULT 0,
    apr numeric,
    due_date timestamp with time zone,
    statement_date timestamp with time zone,
    minimum_payment numeric DEFAULT 0,
    last_payment_date timestamp with time zone,
    last_payment_amount numeric DEFAULT 0,
    rewards_program text,
    active boolean DEFAULT true,
    notes text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.credit_cards OWNER TO daybook_user;

--
-- Name: debt_payments; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.debt_payments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    debt_id uuid NOT NULL,
    account_id uuid NOT NULL,
    amount numeric NOT NULL,
    payment_date timestamp without time zone NOT NULL,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.debt_payments OWNER TO daybook_user;

--
-- Name: debt_records; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.debt_records (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    creditor_name text NOT NULL,
    original_amount numeric NOT NULL,
    remaining_amount numeric NOT NULL,
    account_id uuid,
    status text NOT NULL,
    borrowed_date timestamp without time zone NOT NULL,
    due_date timestamp with time zone,
    interest_rate numeric,
    description text,
    is_initial boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.debt_records OWNER TO daybook_user;

--
-- Name: goal_contributions; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.goal_contributions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    goal_id uuid NOT NULL,
    holding_id uuid,
    type text NOT NULL,
    amount numeric NOT NULL,
    date timestamp with time zone NOT NULL,
    notes text,
    transaction_id uuid NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.goal_contributions OWNER TO daybook_user;

--
-- Name: goal_holdings; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.goal_holdings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    goal_id uuid NOT NULL,
    name text NOT NULL,
    type text NOT NULL,
    status text DEFAULT 'active'::text,
    purchase_date timestamp with time zone NOT NULL,
    amount numeric NOT NULL,
    current_value numeric,
    institution text,
    account_number text,
    interest_rate numeric,
    maturity_date timestamp with time zone,
    maturity_amount numeric,
    tenure_months bigint,
    symbol text,
    quantity numeric,
    cost_basis numeric,
    current_price numeric,
    monthly_deposit numeric,
    details jsonb,
    transaction_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.goal_holdings OWNER TO daybook_user;

--
-- Name: goals; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.goals (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    description text,
    icon text,
    color text,
    category text,
    priority text,
    target_amount numeric NOT NULL,
    current_amount numeric DEFAULT 0,
    target_date timestamp with time zone,
    monthly_contribution numeric,
    status text DEFAULT 'active'::text,
    achieved boolean DEFAULT false,
    achieved_date timestamp with time zone,
    last_contribution numeric,
    last_contribution_date timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.goals OWNER TO daybook_user;

--
-- Name: lend_payments; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.lend_payments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    lend_id uuid NOT NULL,
    account_id uuid NOT NULL,
    amount numeric NOT NULL,
    payment_date timestamp without time zone NOT NULL,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.lend_payments OWNER TO daybook_user;

--
-- Name: lend_records; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.lend_records (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    debtor_name text NOT NULL,
    original_amount numeric NOT NULL,
    remaining_amount numeric NOT NULL,
    account_id uuid,
    status text NOT NULL,
    lent_date timestamp without time zone NOT NULL,
    due_date timestamp with time zone,
    interest_rate numeric,
    description text,
    is_initial boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.lend_records OWNER TO daybook_user;

--
-- Name: reconciliation_transactions; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.reconciliation_transactions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    reconciliation_id uuid NOT NULL,
    transaction_id uuid NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.reconciliation_transactions OWNER TO daybook_user;

--
-- Name: reconciliations; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.reconciliations (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    account_id uuid NOT NULL,
    reconciliation_date timestamp with time zone NOT NULL,
    statement_balance numeric(15,2) NOT NULL,
    book_balance numeric(15,2) NOT NULL,
    difference numeric(15,2) NOT NULL,
    notes text,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.reconciliations OWNER TO daybook_user;

--
-- Name: recurring_transactions; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.recurring_transactions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    template_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    template_user_id uuid NOT NULL,
    template_account_id uuid NOT NULL,
    template_to_account_id uuid,
    template_type text NOT NULL,
    template_amount numeric NOT NULL,
    template_category_id text NOT NULL,
    template_date timestamp with time zone NOT NULL,
    template_description text,
    template_tags jsonb,
    template_savings_goal_id uuid,
    template_fixed_deposit_id uuid,
    template_investment_id uuid,
    template_recurring_id uuid,
    template_credit_card_id uuid,
    template_attachments jsonb,
    template_reconciled boolean DEFAULT false,
    template_reconciliation_id uuid,
    template_created_at timestamp with time zone,
    template_updated_at timestamp with time zone,
    template_deleted_at timestamp with time zone,
    frequency text NOT NULL,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone,
    last_processed timestamp with time zone,
    enabled boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.recurring_transactions OWNER TO daybook_user;

--
-- Name: rewards; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.rewards (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    type text,
    amount numeric,
    description text,
    earned_date timestamp with time zone NOT NULL,
    redeemed boolean DEFAULT false,
    redeemed_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.rewards OWNER TO daybook_user;

--
-- Name: service_records; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.service_records (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    asset_id uuid NOT NULL,
    service_date timestamp without time zone NOT NULL,
    service_type text NOT NULL,
    service_provider text,
    cost numeric NOT NULL,
    description text,
    notes text,
    warranty_covered boolean DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.service_records OWNER TO daybook_user;

--
-- Name: settings; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    currency text DEFAULT 'BDT'::text,
    dark_mode boolean DEFAULT false,
    date_format text DEFAULT 'MM/DD/YYYY'::text,
    first_day_of_week bigint DEFAULT 0,
    language text DEFAULT 'en'::text,
    notif_push boolean,
    notif_email boolean,
    notif_budget_alerts boolean,
    notif_bill_reminders boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.settings OWNER TO daybook_user;

--
-- Name: statements; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.statements (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    statement_date timestamp with time zone NOT NULL,
    due_date timestamp with time zone NOT NULL,
    opening_balance numeric,
    closing_balance numeric,
    minimum_payment numeric,
    total_charges numeric,
    total_payments numeric,
    interest_charged numeric,
    paid boolean DEFAULT false,
    paid_date timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.statements OWNER TO daybook_user;

--
-- Name: tags; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.tags (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    color text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.tags OWNER TO daybook_user;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.transactions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    account_id uuid NOT NULL,
    to_account_id uuid,
    type text NOT NULL,
    amount numeric NOT NULL,
    category_id text NOT NULL,
    date timestamp with time zone NOT NULL,
    description text,
    tags jsonb,
    savings_goal_id uuid,
    fixed_deposit_id uuid,
    investment_id uuid,
    recurring_id uuid,
    credit_card_id uuid,
    attachments jsonb,
    reconciled boolean DEFAULT false,
    reconciliation_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.transactions OWNER TO daybook_user;

--
-- Name: users; Type: TABLE; Schema: public; Owner: daybook_user
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    full_name text,
    role text DEFAULT 'user'::text,
    last_login timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO daybook_user;

--
-- Data for Name: account_types; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.account_types (id, user_id, name, icon, description, active, sort_order, created_at, updated_at, deleted_at) FROM stdin;
d2bdb116-5316-4700-98f7-d788ca98fe78	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Cash	💵	Physical cash	t	1	2025-11-01 11:57:04.989578+00	2025-11-01 11:57:04.989578+00	\N
30614655-80a0-400f-8365-b141d79346d6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Bank	🏦	Bank accounts	t	2	2025-11-01 11:57:04.991716+00	2025-11-01 11:57:04.991716+00	\N
c602b425-44d5-4b44-bab4-f57c61d4cafd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Digital Wallet	📱	Digital payment services	t	3	2025-11-01 11:57:04.993996+00	2025-11-01 11:57:04.993996+00	\N
9c5f9dc6-a64f-45b2-abfd-558f1e2f8f6e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Other	📋	Other account types	t	4	2025-11-01 11:57:04.995559+00	2025-11-01 11:57:04.995559+00	\N
\.


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.accounts (id, user_id, name, type, initial_balance, balance, currency, description, institution, account_number, last_reconciled, reconciliation_difference, active, created_at, updated_at, deleted_at) FROM stdin;
6079ac1c-7a40-45b0-913d-4e2dce11dab6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Brac Bank Debit	bank	355727.12	289202.08	BDT				\N	0	t	2025-11-01 16:59:03.851641+00	2025-11-24 16:03:22.242955+00	\N
338a8a7a-841a-44ca-94c9-1e117406ac70	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	SCB Debit	bank	151650.92	255437.92	BDT				\N	0	t	2025-11-01 12:42:03.199886+00	2025-11-27 16:31:47.312284+00	\N
cfee8a15-8596-439c-9ef3-53bf5252b9f5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Money Bag	cash	50440	65220	BDT				\N	0	t	2025-11-02 17:37:12.975949+00	2025-11-29 03:24:56.762071+00	\N
c3288300-472e-413e-9588-48542b2f66d6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash	digital_wallet	24867.13	3857.9500000000016	BDT				\N	0	t	2025-11-01 12:41:34.973+00	2025-11-29 03:37:05.408492+00	\N
\.


--
-- Data for Name: activity_logs; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.activity_logs (id, user_id, action, module, entity_type, entity_id, description, ip_address, user_agent, metadata, created_at, updated_at, deleted_at) FROM stdin;
0d57f62f-e70d-4170-aaea-c1887ba5f039	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	account	Account	c3288300-472e-413e-9588-48542b2f66d6	Created account: bKash (backfilled)			{"name": "bKash", "type": "digital_wallet", "backfilled": true}	2025-11-01 12:41:34.973+00	2025-11-01 12:41:34.973+00	\N
dd637722-c1ae-4ffe-a9b0-5eec165acfe9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	account	Account	338a8a7a-841a-44ca-94c9-1e117406ac70	Created account: SCB Debit (backfilled)			{"name": "SCB Debit", "type": "bank", "backfilled": true}	2025-11-01 12:42:03.199886+00	2025-11-01 12:42:03.199886+00	\N
bc0c8d5e-fd0f-4e26-8625-46a84887c618	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	account	Account	6079ac1c-7a40-45b0-913d-4e2dce11dab6	Created account: Brac Bank Debit (backfilled)			{"name": "Brac Bank Debit", "type": "bank", "backfilled": true}	2025-11-01 16:59:03.851641+00	2025-11-01 16:59:03.851641+00	\N
4154ba75-54ec-41c3-974f-5645a8197535	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	account	Account	cfee8a15-8596-439c-9ef3-53bf5252b9f5	Created account: Money Bag (backfilled)			{"name": "Money Bag", "type": "cash", "backfilled": true}	2025-11-02 17:37:12.975949+00	2025-11-02 17:37:12.975949+00	\N
6ba5331b-7c54-45ca-bf94-374d8978c5e2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	20927be9-b4f8-4c1a-bd25-4a61bc1fbb8a	Created transaction: Opening balance for bKash (backfilled)			{"type": "income", "amount": 24867.13, "backfilled": true}	2025-11-01 12:41:34.977901+00	2025-11-01 12:41:34.977901+00	\N
da10b401-9faa-4d77-86e2-70d2dc58a1a2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	2358ced5-3577-4d62-9fd6-015de3d01686	Created transaction: Opening balance for SCB Debit (backfilled)			{"type": "income", "amount": 151650.92, "backfilled": true}	2025-11-01 12:42:03.200322+00	2025-11-01 12:42:03.200322+00	\N
c9ebc0cf-4a9d-4dfa-8ee2-553aa0788bbb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	68a47ac4-8087-47b1-ae4c-410628e07c2f	Created transaction: Opening balance for Brac Bank Debit (backfilled)			{"type": "income", "amount": 355727.12, "backfilled": true}	2025-11-01 16:59:03.852271+00	2025-11-01 16:59:03.852271+00	\N
69d91253-e0bb-4585-b09d-feccf0d96512	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	05000557-7353-4e92-8c70-971e8912a77c	Created transaction: External Fixed Deposit tracked for Brac Bank Fixed Deposit (backfilled)			{"type": "tracking", "amount": 200000, "backfilled": true}	2025-11-01 18:17:53.864547+00	2025-11-01 18:17:53.864547+00	\N
6822db77-672a-4d7a-9825-c9fe512cddde	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1b9f8926-a5ba-4308-aa5c-45917b2149f7	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 112, "backfilled": true}	2025-11-02 15:48:34.224044+00	2025-11-02 15:48:34.224044+00	\N
7fc0257f-33bf-4dd4-8cc4-bc63f2b6a029	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c3802ca2-dbc0-4652-8065-7782a3b2b1e9	Created transaction: Donation for Nuts (backfilled)			{"type": "expense", "amount": 499, "backfilled": true}	2025-11-02 15:49:50.059598+00	2025-11-02 15:49:50.059598+00	\N
e81880cb-6210-43fd-ab81-078c9decb1b9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	bba48dc4-2b64-4a1d-9bed-4e8b9ffeaaf3	Created transaction: External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230866410 (backfilled)			{"type": "tracking", "amount": 4000, "backfilled": true}	2025-11-02 16:08:00.233165+00	2025-11-02 16:08:00.233165+00	\N
5da47b04-ec54-4635-b9c0-6fece1c222ae	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1f91f320-bac8-4739-897b-0dfeed161183	Created transaction: External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230701614 (backfilled)			{"type": "tracking", "amount": 50000, "backfilled": true}	2025-11-02 16:21:46.399021+00	2025-11-02 16:21:46.399021+00	\N
af54266a-d671-46bf-a6ea-c2bbbe8bb203	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	fe8cf6ae-5ead-4ccf-a551-10289a99324b	Created transaction: External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230291760 (backfilled)			{"type": "tracking", "amount": 10000, "backfilled": true}	2025-11-02 16:24:02.665361+00	2025-11-02 16:24:02.665361+00	\N
27b7e7cc-53bb-42de-8c62-165002037279	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1cbdee6a-62bd-48b2-bf4d-b99f2984c1f6	Created transaction: External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230221980 (backfilled)			{"type": "tracking", "amount": 10000, "backfilled": true}	2025-11-02 16:35:03.543314+00	2025-11-02 16:35:03.543314+00	\N
cb34e698-8267-47a1-962f-646936db5413	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e62ea004-113d-4508-9ec2-21ee2c926bd5	Created transaction: External DPS (Deposit Pension Scheme) tracked for bKash DPS - 1783060326999 (backfilled)			{"type": "tracking", "amount": 13000, "backfilled": true}	2025-11-02 16:38:02.101328+00	2025-11-02 16:38:02.101328+00	\N
44771e88-c521-4a0f-b4d3-c475f4a32777	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	6b612542-991f-40c7-94d5-2c0ab3756034	Created transaction: External DPS (Deposit Pension Scheme) tracked for bKash DPS - 1783060302607 (backfilled)			{"type": "tracking", "amount": 55000, "backfilled": true}	2025-11-02 16:43:44.831818+00	2025-11-02 16:43:44.831818+00	\N
d5402197-2308-4d16-813c-eb11ccd06831	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	faee24b9-b925-48a0-87f4-268541d59b03	Created transaction: Added to bKash DPS - 2145230866410: DPS (Deposit Pension Scheme) (backfilled)			{"type": "expense", "amount": 1000, "backfilled": true}	2025-11-02 16:47:01.998419+00	2025-11-02 16:47:01.998419+00	\N
22fdfa98-d5da-420e-adf8-ab8e5308956a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c6ee0331-3670-4842-a99b-1a5cbe41b92f	Created transaction: Added to bKash DPS - 2145230701614: DPS (Deposit Pension Scheme) (backfilled)			{"type": "expense", "amount": 10000, "backfilled": true}	2025-11-02 16:47:33.517648+00	2025-11-02 16:47:33.517648+00	\N
e80727ae-adcc-4411-9854-4aa3f54c5857	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	2e95fe7b-1114-4a03-bb4e-c3c2d526b741	Created transaction: Added to bKash DPS - 2145230291760: DPS (Deposit Pension Scheme) (backfilled)			{"type": "expense", "amount": 10000, "backfilled": true}	2025-11-02 16:48:03.748527+00	2025-11-02 16:48:03.748527+00	\N
fc41cca5-8096-4022-af73-f450ab8a9be2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	98ab9145-1d9b-4258-991c-4ec6c145ea68	Created transaction: Added to bKash DPS - 2145230221980: DPS (Deposit Pension Scheme) (backfilled)			{"type": "expense", "amount": 1000, "backfilled": true}	2025-11-02 16:50:33.529188+00	2025-11-02 16:50:33.529188+00	\N
fd53b2c6-9fc9-49d8-a036-26119a9686b3	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a71302a1-0030-4284-b654-278088dd54bf	Created transaction: Internet bill (backfilled)			{"type": "expense", "amount": 1500, "backfilled": true}	2025-11-02 16:51:05.285357+00	2025-11-02 16:51:05.285357+00	\N
55370bcb-ca97-4b00-ac3c-a6160fc6831c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	fad607ef-b4f7-415b-804f-d1655b0fa695	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 112, "backfilled": true}	2025-11-02 16:51:44.236531+00	2025-11-02 16:51:44.236531+00	\N
0fb9564d-5358-4eeb-8c93-e9747c4f87aa	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	96f5ce51-59e5-4257-af1d-40f45559b1e7	Created transaction: DESCO prepaid bill (backfilled)			{"type": "expense", "amount": 500, "backfilled": true}	2025-11-02 16:52:39.859276+00	2025-11-02 16:52:39.859276+00	\N
62fc9ae6-a456-4ab4-ac30-bc3db95b2d6d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	5ba2b13e-0ec4-4caf-a755-8e716f9a3068	Created transaction: External Fixed Deposit tracked for IDLC ISF-SIP-001449 (backfilled)			{"type": "tracking", "amount": 725000, "backfilled": true}	2025-11-02 17:10:37.264537+00	2025-11-02 17:10:37.264537+00	\N
757ec6ae-ac76-4b1d-9b91-2f082fd39527	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	2d74ddb2-48a7-493a-bbf8-caf8d029f251	Created transaction: Added to IDLC ISF-SIP-001449: Fixed Deposit (backfilled)			{"type": "expense", "amount": 25000, "backfilled": true}	2025-11-02 17:11:10.005217+00	2025-11-02 17:11:10.005217+00	\N
6f2a2c89-3611-4a71-8f58-3bde34fac2f3	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	24e09bac-0410-4256-b9ff-2136ed55feff	Created transaction: External Fixed Deposit tracked for IDLC ISF-SIP-001807 (backfilled)			{"type": "tracking", "amount": 392000, "backfilled": true}	2025-11-02 17:16:55.899574+00	2025-11-02 17:16:55.899574+00	\N
c42cb5f7-3141-42f8-8435-76c6c5c6b483	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c4ca3200-5561-448f-ad23-56546fe3fc52	Created transaction: Added to IDLC ISF-SIP-001807: Fixed Deposit (backfilled)			{"type": "expense", "amount": 12000, "backfilled": true}	2025-11-02 17:17:16.924063+00	2025-11-02 17:17:16.924063+00	\N
c3c28b6f-a224-44dc-9245-6b934f4d1c9f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	ef95c42e-a4bf-4330-9fc9-431bd67f94e8	Created transaction: Opening balance for Money Bag (backfilled)			{"type": "income", "amount": 50440, "backfilled": true}	2025-11-02 17:37:12.978036+00	2025-11-02 17:37:12.978036+00	\N
1fe54751-db58-49d3-bdcd-ef015b1e8524	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c099b99b-3e6b-480c-8d5e-da9eb84a77b6	Created transaction: Bike fare (backfilled)			{"type": "expense", "amount": 140, "backfilled": true}	2025-11-02 17:37:31.666623+00	2025-11-02 17:37:31.666623+00	\N
6baee0ad-1a0e-4c67-82f4-17061a04b02e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a646668c-1f77-408e-bff8-165e480a4d89	Created transaction: Rickshaw fare (backfilled)			{"type": "expense", "amount": 100, "backfilled": true}	2025-11-02 17:37:52.604777+00	2025-11-02 17:37:52.604777+00	\N
b3807dcc-e4dc-4a92-bdc0-3805fd3b894d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c702dfba-6359-4b4f-9777-8e67e53c7829	Created transaction: Rickshaw fare (backfilled)			{"type": "expense", "amount": 100, "backfilled": true}	2025-11-03 05:20:14.088709+00	2025-11-03 05:20:14.088709+00	\N
6be94956-6fd2-4fff-8fb4-3a08923ea054	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	abbafb8a-4521-4155-815e-82af3798f4e9	Created transaction: Garbage Bag for Mess (backfilled)			{"type": "expense", "amount": 130, "backfilled": true}	2025-11-03 15:53:42.270723+00	2025-11-03 15:53:42.270723+00	\N
a6c2ea01-ee0c-4d94-845d-e9ab4294d57f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	eb273fef-13dd-419d-92f7-c425b6d0d721	Created transaction: Rickshaw fare (backfilled)			{"type": "expense", "amount": 100, "backfilled": true}	2025-11-03 15:54:01.348217+00	2025-11-03 15:54:01.348217+00	\N
d252818e-d0fe-4a93-918b-42767d52b382	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	acb6339d-1d2c-4257-9860-2c38264b7a4f	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 112, "backfilled": true}	2025-11-03 15:54:40.473078+00	2025-11-03 15:54:40.473078+00	\N
ce453907-d8db-4144-99ba-2b4d3ff2d178	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3c02d554-581e-4c36-8b3a-f1bb37709218	Created transaction: Transfer between accounts (backfilled)			{"type": "transfer", "amount": 5000, "backfilled": true}	2025-11-04 13:31:16.59+00	2025-11-04 13:31:16.59+00	\N
6598367a-1a7c-4d82-b2bc-47a66539da30	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	065c91fe-30c3-4e06-abdc-3c0a76733a52	Created transaction: Secret recipe (backfilled)			{"type": "expense", "amount": 1505, "backfilled": true}	2025-11-05 04:09:14.787251+00	2025-11-05 04:09:14.787251+00	\N
cdf420e0-94f7-4afd-9ad6-850949339773	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	97baf158-9df0-4df0-8936-4e237c2a3250	Created transaction: Rickshaw fare (backfilled)			{"type": "expense", "amount": 100, "backfilled": true}	2025-11-05 04:09:35.761072+00	2025-11-05 04:09:35.761072+00	\N
10b0ad2f-6cc3-410e-89f6-efb632825559	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	6a66ed19-e95a-4e2e-b349-6be2125f5a5e	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 10, "backfilled": true}	2025-11-05 04:09:58.578124+00	2025-11-05 04:09:58.578124+00	\N
ceda8983-02bb-4e8a-9c96-f94373e16d8b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	eb8be643-5247-4672-a930-b6dd813cddf6	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 200, "backfilled": true}	2025-11-05 04:10:46.330631+00	2025-11-05 04:10:46.330631+00	\N
28158e63-58a9-42d2-b0ba-58880a7d4cdf	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	5cad292a-0597-4ea2-976e-200ec3af2df7	Created transaction: Rickshaw fare (backfilled)			{"type": "expense", "amount": 100, "backfilled": true}	2025-11-05 15:52:58.350899+00	2025-11-05 15:52:58.350899+00	\N
a58ac99a-4e06-4a29-8913-d9164e7be6f5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3affaafc-ec1d-4487-b823-d9315edc8007	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 112, "backfilled": true}	2025-11-05 15:54:50.720092+00	2025-11-05 15:54:50.720092+00	\N
0f027d6c-a105-4407-b930-4fda0759cf04	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	0df44019-416c-4654-bdb8-ad5368f2cdc7	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 112, "backfilled": true}	2025-11-07 08:42:47.814545+00	2025-11-07 08:42:47.814545+00	\N
0478ffdf-3440-4caa-ab7b-432b896411c7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e4976e8d-add7-4c47-9cfe-14d6a5c988e1	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 112, "backfilled": true}	2025-11-07 08:43:01.602677+00	2025-11-07 08:43:01.602677+00	\N
540f8948-6cb7-4def-8ab0-4275cfeb37bf	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c5f5e973-f8d7-40c6-a473-304ba50e8e36	Created transaction: Bike fare (backfilled)			{"type": "expense", "amount": 90, "backfilled": true}	2025-11-07 08:43:16.878293+00	2025-11-07 08:43:16.878293+00	\N
7d491524-e6b7-4c8d-9f78-abc809e09b54	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	5791f38a-7e9a-4f11-a4d9-052b64b59c54	Created transaction: Sarah resort pepsi (backfilled)			{"type": "expense", "amount": 400, "backfilled": true}	2025-11-07 08:43:44.962968+00	2025-11-07 08:43:44.962968+00	\N
25e6bcee-b32b-475b-8ea5-eb959adf589f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7f51ae39-3f79-4d8e-9ba9-7e1c1436e0ab	Created transaction: Rickshaw fare (backfilled)			{"type": "expense", "amount": 70, "backfilled": true}	2025-11-07 13:08:27.5759+00	2025-11-07 13:08:27.5759+00	\N
f36f5201-ce1f-49b7-9d87-0e26dc095f0e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	f04b6cf5-49cf-4c3b-a20d-9a2fa8a94dad	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 112, "backfilled": true}	2025-11-08 06:46:44.986574+00	2025-11-08 06:46:44.986574+00	\N
b99bb99b-d1a1-474f-9e1c-7dcb5432ccf5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	budget	Budget	f6f03cd1-a6a0-4406-8a0b-dc8f98935c15	Created budget for category transport (backfilled)			{"amount": 5000, "period": "monthly", "backfilled": true, "categoryId": "transport"}	2025-11-07 13:05:54.604167+00	2025-11-07 13:05:54.604167+00	\N
1129b964-a6c6-4812-b5d4-57e41ad6af79	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	budget	Budget	6754bb1e-d764-43c6-8e63-9347dcfbbfec	Created budget for category food (backfilled)			{"amount": 10000, "period": "monthly", "backfilled": true, "categoryId": "food"}	2025-11-07 13:06:03.740615+00	2025-11-07 13:06:03.740615+00	\N
25d45afe-bda0-4469-b6dd-39e40ecee340	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	budget	Budget	07df20f4-a407-4a7f-bada-918c5d04d3fb	Created budget for category donation (backfilled)			{"amount": 5000, "period": "monthly", "backfilled": true, "categoryId": "donation"}	2025-11-07 13:06:13.974766+00	2025-11-07 13:06:13.974766+00	\N
8f48f39f-7d7b-4b17-864e-8dbcff3a956b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	credit_card	CreditCard	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	Created credit card: Brac Bank (Credit Card) (backfilled)			{"name": "Brac Bank (Credit Card)", "backfilled": true}	2025-11-01 17:00:54.255539+00	2025-11-01 17:00:54.255539+00	\N
84bae6b2-6f08-40d4-af69-207fb41df347	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	lend	LendRecord	0f34ad12-adea-40f5-ac80-0cde89b0c630	Created lend record: Sohug Mama (backfilled)			{"amount": 140000, "backfilled": true}	2025-11-08 08:25:09.904137+00	2025-11-08 08:25:09.904137+00	\N
f3cf1acf-01ff-4b06-a4a6-2a5e9e2dbd22	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	dd8540c6-c548-4989-aa0f-41b1b379bafb	Created goal: IDLC ISF-SIP-001449 (backfilled)			{"name": "IDLC ISF-SIP-001449", "backfilled": true, "target_amount": 1500000}	2025-11-02 17:09:18.477154+00	2025-11-02 17:09:18.477154+00	\N
1583c8bc-4044-47da-9474-8aa3345883af	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	6ab6926d-8ede-41c0-8cb0-a2f36ce6ea11	Created goal: Brac Bank Fixed Deposit (backfilled)			{"name": "Brac Bank Fixed Deposit", "backfilled": true, "target_amount": 200000}	2025-11-01 18:16:49.105765+00	2025-11-01 18:16:49.105765+00	\N
b7a0aa6c-8de8-41fa-9d26-f8922d43da53	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	df2b69ff-147f-4676-af50-25db9ff1ee54	Created goal: bKash DPS - 2145230866410 (backfilled)			{"name": "bKash DPS - 2145230866410", "backfilled": true, "target_amount": 12000}	2025-11-02 16:06:52.837283+00	2025-11-02 16:06:52.837283+00	\N
832d4f62-4074-4335-be87-bfe172b32e0e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	e0dc0db0-05c6-4ace-9f5d-6755859463db	Created goal: bKash DPS - 2145230701614 (backfilled)			{"name": "bKash DPS - 2145230701614", "backfilled": true, "target_amount": 120000}	2025-11-02 16:20:40.713877+00	2025-11-02 16:20:40.713877+00	\N
83509f6a-c551-48cd-be3f-9f1646a14c38	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	e6b69c0c-d123-4f9b-a773-c818a5dfdaa9	Created goal: bKash DPS - 2145230291760 (backfilled)			{"name": "bKash DPS - 2145230291760", "backfilled": true, "target_amount": 120000}	2025-11-02 16:23:03.244867+00	2025-11-02 16:23:03.244867+00	\N
a49c7078-82e9-4515-b26b-0403eebae654	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	84eb448d-76dd-46c6-a4b7-9c9179f6753e	Created goal: bKash DPS - 2145230221980 (backfilled)			{"name": "bKash DPS - 2145230221980", "backfilled": true, "target_amount": 12000}	2025-11-02 16:34:21.911244+00	2025-11-02 16:34:21.911244+00	\N
9da4ba90-026c-444e-b012-4c971a39609b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	56ff16a0-de40-4241-8904-f6680afb2a61	Created goal: bKash DPS - 1783060326999 (backfilled)			{"name": "bKash DPS - 1783060326999", "backfilled": true, "target_amount": 24000}	2025-11-02 16:37:17.7819+00	2025-11-02 16:37:17.7819+00	\N
4c089ac0-1b5e-447b-8ebb-c1a9575d148e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	33e31fd1-12ae-4df7-b31c-fda997d4e4ce	Created goal: bKash DPS - 1783060302607 (backfilled)			{"name": "bKash DPS - 1783060302607", "backfilled": true, "target_amount": 72000}	2025-11-02 16:43:04.930308+00	2025-11-02 16:43:04.930308+00	\N
0a0a71d4-fbd7-45cb-b422-26a068e6463d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	Goal	2946b607-6a79-4617-8c2b-bbebd47657e5	Created goal: IDLC ISF-SIP-001807 (backfilled)			{"name": "IDLC ISF-SIP-001807", "backfilled": true, "target_amount": 720000}	2025-11-02 17:14:11.019772+00	2025-11-02 17:14:11.019772+00	\N
8599764c-2cf8-4c22-88ff-a522083c5725	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	a39f51d8-52a8-4fd5-b2dd-6e01f27fe0a9	Created asset: GRID Monitor Stand 	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 14:46:35.653412+00	2025-11-08 14:46:35.653412+00	\N
373abaec-b553-41a5-a9f7-565daf1f36a6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3d1b594e-f394-4600-9406-253e43b14621	Created transaction: Mobile recharge for Maa	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 15:06:28.166173+00	2025-11-08 15:06:28.166173+00	\N
ba915895-1719-4a7d-9d9d-31d56af8238e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	3773b103-2b4e-4a97-92bc-02711c7dfe3d	Created asset: GRID Industrious Work Desk	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 15:19:09.013862+00	2025-11-08 15:19:09.013862+00	\N
8516f69d-44ed-4aa0-8856-3631b88062fb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	8fce39d8-3365-429c-a95e-f4c220d92e30	Created asset: GRID Filing Cabinet	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 15:21:39.48209+00	2025-11-08 15:21:39.48209+00	\N
71b79f37-21c9-4691-8e79-6217fbcbf475	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	538e0600-4050-45c7-bd63-0193baabaf47	Created asset: GRID Comfy Chair	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 15:23:29.935919+00	2025-11-08 15:23:29.935919+00	\N
70e1178d-ae7a-4740-bf93-9642cd9b48b7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	c701367f-53d9-4354-96a3-dc982e0eca9a	Created asset: Seagate STHH2000400 Backup plus Ultra touch 2TB USB-C Portable Hard Drive	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 15:33:25.894813+00	2025-11-08 15:33:25.894813+00	\N
2c512feb-5e27-44e0-9366-7c4c6d5de392	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	ed6a60bc-86f8-426d-af91-466e70edf585	Created asset: Ugreen PB760 10000mAh 20W Magnetic Wireless Power Bank #35341	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 16:00:05.720554+00	2025-11-08 16:00:05.720554+00	\N
5aaa401b-2116-4da4-8ff0-f4144a71b056	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	c8256dd6-c1fa-4fd0-959b-64d3fec1d2fe	Created asset: Havit H94 4-Port High-Speed USB Hub	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 16:02:10.071481+00	2025-11-08 16:02:10.071481+00	\N
3acade04-597d-4b0f-adf4-116a038852e1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	928306b1-1bae-42c7-bab7-4c7f9a6b4569	Created asset: UGREEN 60126 USB-A 2.0 To Type-C 1 Meter Data Cable	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 16:03:18.029368+00	2025-11-08 16:03:18.029368+00	\N
d410656a-35e5-44d0-baa6-6a31be8aa445	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	6ce7eda6-068a-4074-a33e-2047d589d9fe	Created asset: Transcend 64GB Micro SD UHS-I U1 Memory Card with Adapter (TS64GUSD300S-A)	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 16:05:19.591905+00	2025-11-08 16:05:19.591905+00	\N
430d6509-7366-4ad3-bfc7-4493db1e8515	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	e78b4635-1c7d-4f38-b76d-947ae97ee11d	Created asset: UiiSii HM13 Wired In-Ear Headphone with Mic	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 16:14:36.625785+00	2025-11-08 16:14:36.625785+00	\N
66cb5240-926b-4e11-a907-c31dc90da802	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	7725934c-624b-413d-adac-40bde09e1cc1	Created asset: TP-Link Deco M4 (Single Pack) Whole Home Mesh Wi-Fi System AC1200 Dual-band Router	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 16:33:04.544893+00	2025-11-08 16:33:04.544893+00	\N
ef9f13fd-6d49-4652-ab42-cfd0b1839a33	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	1ffb3514-d3b2-4674-9033-cfe29856d983	Created asset: Apple AirPods Pro	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 16:38:31.90766+00	2025-11-08 16:38:31.90766+00	\N
b9527afb-04b3-4e62-bf6c-2a87dd67acfe	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	d62f9d57-7874-4292-8497-d137e3a8a609	Created transaction: Guava and Rice bran oil	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:02:07.648944+00	2025-11-08 18:02:07.648944+00	\N
3963f034-0590-440e-97a5-33305f689cf7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1ee7cc20-fa4d-4693-90c0-8ccd1cff5ce2	Created transaction: Food for Mess	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:02:33.190241+00	2025-11-08 18:02:33.190241+00	\N
a9801012-2813-414d-bfb0-b80558e6fbfd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	1ee7cc20-fa4d-4693-90c0-8ccd1cff5ce2	Updated transaction: Food for Mess	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:02:56.944328+00	2025-11-08 18:02:56.944328+00	\N
fbffef6e-a57e-4feb-ad4b-58f4b67ae2a1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	1ee7cc20-fa4d-4693-90c0-8ccd1cff5ce2	Updated transaction: Food for Mess	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:03:29.008369+00	2025-11-08 18:03:29.008369+00	\N
9ecce5dd-9301-491f-b630-153f938e73b4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	d62f9d57-7874-4292-8497-d137e3a8a609	Updated transaction: Guava and Rice bran oil	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:04:20.414024+00	2025-11-08 18:04:20.414024+00	\N
48b6c3cb-7b63-41bc-9b4e-0363aedd905d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	d62f9d57-7874-4292-8497-d137e3a8a609	Updated transaction: Guava and Rice bran oil	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:04:33.57094+00	2025-11-08 18:04:33.57094+00	\N
4ec1e756-18ad-424a-8ec1-2d04a7fe6776	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	credit_card	CreditCardPayment	f07125fd-9bd2-43e7-ba3a-d1d2af141cab	Recorded credit card payment: Brac Bank (Credit Card)	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:04:44.464065+00	2025-11-08 18:04:44.464065+00	\N
070575b8-faa1-4c1d-8bf5-809ec4da3efc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	lend	Lend	55821a5e-b9f0-4cb4-a60f-084f5bd76f44	Created lend record for Ripon Bhai	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-08 18:07:36.392759+00	2025-11-08 18:07:36.392759+00	\N
03bdd831-1876-476d-90ce-9c4de98b6192	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3d1b594e-f394-4600-9406-253e43b14621	Created transaction: Mobile recharge for Maa (backfilled)			{"type": "expense", "amount": 598, "backfilled": true}	2025-11-08 15:06:28.161034+00	2025-11-08 15:06:28.161034+00	\N
ade7fbfe-496c-4612-85ff-451dffc60975	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1ee7cc20-fa4d-4693-90c0-8ccd1cff5ce2	Created transaction: Food for Mess (backfilled)			{"type": "expense", "amount": 1014, "backfilled": true}	2025-11-08 18:02:33.184802+00	2025-11-08 18:02:33.184802+00	\N
6024f738-0c42-4597-b32e-0f7b6e48f94a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	d62f9d57-7874-4292-8497-d137e3a8a609	Created transaction: Guava and Rice bran oil (backfilled)			{"type": "expense", "amount": 600.81, "backfilled": true}	2025-11-08 18:02:07.645428+00	2025-11-08 18:02:07.645428+00	\N
b6fe0321-2a14-42d9-80f3-bc154b4a6e63	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	ad1ac008-a63b-4ecc-999b-bbe1ad51e3c2	Created transaction: Credit card payment - Brac Bank (Credit Card) (backfilled)			{"type": "expense", "amount": 1614.81, "backfilled": true}	2025-11-08 18:04:44.455194+00	2025-11-08 18:04:44.455194+00	\N
880fd52e-ba75-4427-9df8-fd8804464024	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	lend	LendRecord	55821a5e-b9f0-4cb4-a60f-084f5bd76f44	Created lend record: Ripon Bhai (backfilled)			{"amount": 277500, "backfilled": true}	2025-11-08 18:07:36.39057+00	2025-11-08 18:07:36.39057+00	\N
049f6a60-6f5b-489d-8fc5-cca3f573a317	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	a39f51d8-52a8-4fd5-b2dd-6e01f27fe0a9	Created asset: GRID Monitor Stand  (backfilled)			{"name": "GRID Monitor Stand ", "backfilled": true, "purchasePrice": 2500}	2025-11-08 14:46:35.64897+00	2025-11-08 14:46:35.64897+00	\N
ad6c722d-0666-40c1-8b68-758f88cabf6a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	3773b103-2b4e-4a97-92bc-02711c7dfe3d	Created asset: GRID Industrious Work Desk (backfilled)			{"name": "GRID Industrious Work Desk", "backfilled": true, "purchasePrice": 5000}	2025-11-08 15:19:09.007209+00	2025-11-08 15:19:09.007209+00	\N
a05d9759-f57c-42f4-885c-5a4ef1ee17e2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	8fce39d8-3365-429c-a95e-f4c220d92e30	Created asset: GRID Filing Cabinet (backfilled)			{"name": "GRID Filing Cabinet", "backfilled": true, "purchasePrice": 5500}	2025-11-08 15:21:39.477664+00	2025-11-08 15:21:39.477664+00	\N
a0b4127a-d539-49c1-ac94-4b9f926af78e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	538e0600-4050-45c7-bd63-0193baabaf47	Created asset: GRID Comfy Chair (backfilled)			{"name": "GRID Comfy Chair", "backfilled": true, "purchasePrice": 16500}	2025-11-08 15:23:29.93397+00	2025-11-08 15:23:29.93397+00	\N
b0cb49f1-fb5c-4b21-aa28-d622f3cb6d0f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	c701367f-53d9-4354-96a3-dc982e0eca9a	Created asset: Seagate STHH2000400 Backup plus Ultra touch 2TB USB-C Portable Hard Drive (backfilled)			{"name": "Seagate STHH2000400 Backup plus Ultra touch 2TB USB-C Portable Hard Drive", "backfilled": true, "purchasePrice": 9000}	2025-11-08 15:33:25.891351+00	2025-11-08 15:33:25.891351+00	\N
82f695d4-5a41-4cfc-befe-c396cf2ffd27	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	ed6a60bc-86f8-426d-af91-466e70edf585	Created asset: Ugreen PB760 10000mAh 20W Magnetic Wireless Power Bank #35341 (backfilled)			{"name": "Ugreen PB760 10000mAh 20W Magnetic Wireless Power Bank #35341", "backfilled": true, "purchasePrice": 4020}	2025-11-08 16:00:05.718367+00	2025-11-08 16:00:05.718367+00	\N
1e7bb10c-998e-4171-944d-435099d46831	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	c8256dd6-c1fa-4fd0-959b-64d3fec1d2fe	Created asset: Havit H94 4-Port High-Speed USB Hub (backfilled)			{"name": "Havit H94 4-Port High-Speed USB Hub", "backfilled": true, "purchasePrice": 1025}	2025-11-08 16:02:10.069017+00	2025-11-08 16:02:10.069017+00	\N
e8afaf05-2c07-488c-b0e5-4e3e288481ff	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	928306b1-1bae-42c7-bab7-4c7f9a6b4569	Created asset: UGREEN 60126 USB-A 2.0 To Type-C 1 Meter Data Cable (backfilled)			{"name": "UGREEN 60126 USB-A 2.0 To Type-C 1 Meter Data Cable", "backfilled": true, "purchasePrice": 325}	2025-11-08 16:03:18.026248+00	2025-11-08 16:03:18.026248+00	\N
636dd71b-780e-4277-a35a-a68638ec6942	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	6ce7eda6-068a-4074-a33e-2047d589d9fe	Created asset: Transcend 64GB Micro SD UHS-I U1 Memory Card with Adapter (TS64GUSD300S-A) (backfilled)			{"name": "Transcend 64GB Micro SD UHS-I U1 Memory Card with Adapter (TS64GUSD300S-A)", "backfilled": true, "purchasePrice": 800}	2025-11-08 16:05:19.589561+00	2025-11-08 16:05:19.589561+00	\N
4d387983-7652-41c5-b3cb-d9d63d7fe196	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	e78b4635-1c7d-4f38-b76d-947ae97ee11d	Created asset: UiiSii HM13 Wired In-Ear Headphone with Mic (backfilled)			{"name": "UiiSii HM13 Wired In-Ear Headphone with Mic", "backfilled": true, "purchasePrice": 730}	2025-11-08 16:14:36.62371+00	2025-11-08 16:14:36.62371+00	\N
183e3bb9-01b8-4c77-894d-8db21428ba28	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	7725934c-624b-413d-adac-40bde09e1cc1	Created asset: TP-Link Deco M4 (Single Pack) Whole Home Mesh Wi-Fi System AC1200 Dual-band Router (backfilled)			{"name": "TP-Link Deco M4 (Single Pack) Whole Home Mesh Wi-Fi System AC1200 Dual-band Router", "backfilled": true, "purchasePrice": 5200}	2025-11-08 16:33:04.542262+00	2025-11-08 16:33:04.542262+00	\N
92d0fda7-05ee-4bff-9f09-9b3ed4ca886f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	1ffb3514-d3b2-4674-9033-cfe29856d983	Created asset: Apple AirPods Pro (backfilled)			{"name": "Apple AirPods Pro", "backfilled": true, "purchasePrice": 21990}	2025-11-08 16:38:31.90546+00	2025-11-08 16:38:31.90546+00	\N
d885f198-0509-4e70-acf7-5491c407206e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	af3f64b8-716c-4a1a-85da-d67543d3e588	Created transaction: Bua Bill	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 16:57:42.339619+00	2025-11-09 16:57:42.339619+00	\N
5508c674-cae5-4d6e-b02a-afd9e1f1126c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	af3f64b8-716c-4a1a-85da-d67543d3e588	Updated transaction: Bua Bill	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 17:52:57.219384+00	2025-11-09 17:52:57.219384+00	\N
13032656-325b-49cd-b3dc-a485375454c1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	be686fc0-a2a1-40b2-bb58-adb7d4e85e65	Created asset: Apple MacBook Pro M1	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 18:04:50.265445+00	2025-11-09 18:04:50.265445+00	\N
8a2ceb0d-448d-49ef-9f14-f1abc44f5928	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	asset	Asset	69509677-7d3f-4a00-83a6-509d013171b0	Updated asset: Raspberry Pi 4	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 18:09:45.42566+00	2025-11-09 18:09:45.42566+00	\N
3b99a217-8e3a-41c5-9d1a-a9458aa37682	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:09:18.31496+00	2025-11-10 08:09:18.31496+00	\N
72659ae0-6974-40f7-8c64-8bdf943f0e58	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:09:21.026035+00	2025-11-10 08:09:21.026035+00	\N
350e9d7d-bcc0-4b4f-a995-47c874b21b7c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:09:23.847075+00	2025-11-10 08:09:23.847075+00	\N
87f0212e-ee13-4243-b2dd-eed37d834b3a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:09:25.0453+00	2025-11-10 08:09:25.0453+00	\N
2453102d-195d-4c1f-bb6d-0cf87962a5ad	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:09:25.954774+00	2025-11-10 08:09:25.954774+00	\N
8e8e6057-cfcc-42b9-9dad-3f279261e66a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:09:26.806407+00	2025-11-10 08:09:26.806407+00	\N
e1e48563-0b73-4289-88ca-a31250191c03	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:21.758976+00	2025-11-10 08:11:21.758976+00	\N
7a39fe0e-671c-48db-8e3f-729070647ba9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:23.568069+00	2025-11-10 08:11:23.568069+00	\N
0354b580-64a1-4f0b-b0fd-7304cd3ebce4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:25.771313+00	2025-11-10 08:11:25.771313+00	\N
dee1eaca-d335-450e-ad76-afa4f5c6fd5c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 16:57:11.923419+00	2025-11-09 16:57:11.923419+00	\N
a40899c9-749f-41bd-aec0-6ae9d6347880	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c272af47-ed23-473a-b866-1ef1269d2aaa	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 16:58:14.78463+00	2025-11-09 16:58:14.78463+00	\N
245c436e-298a-4497-8503-c3147e553d32	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	9e4fad59-4d4b-4bf6-97d5-dfc1200ec81f	Created transaction: CNG fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 16:58:35.864653+00	2025-11-09 16:58:35.864653+00	\N
4d60a544-05d9-4f9e-b1eb-5b68d2fa0a4c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	87b3d61e-80ca-4e14-99fc-1d0763d98efd	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 16:58:53.66108+00	2025-11-09 16:58:53.66108+00	\N
3d183856-86e8-472b-848f-0b14d3244bf6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1394b9c5-a5bb-4d2d-b265-610a264a4c20	Created transaction: TAX pay	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 16:59:45.645188+00	2025-11-09 16:59:45.645188+00	\N
8011faf3-3171-4474-a8cc-7915a9ec3dc2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e75f993c-d932-47f1-b47a-61d116364867	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 17:01:16.605436+00	2025-11-09 17:01:16.605436+00	\N
8ff998bb-f954-49d0-b79a-44f1db56d489	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	69509677-7d3f-4a00-83a6-509d013171b0	Created asset: Raspberry Pi 4	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-09 17:56:39.297539+00	2025-11-09 17:56:39.297539+00	\N
b51bdf30-d431-4058-9fa0-e51b2475d110	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:26.636753+00	2025-11-10 08:11:26.636753+00	\N
c23df5a6-c44f-423c-a486-10d971a4c8d2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:27.140078+00	2025-11-10 08:11:27.140078+00	\N
85c433d4-c27b-4dc6-b3fd-9a48fd48da85	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:27.439117+00	2025-11-10 08:11:27.439117+00	\N
0724c00d-49ba-41d3-a2f1-fe8bd7cc67ea	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:27.790842+00	2025-11-10 08:11:27.790842+00	\N
3b77d1f1-75cd-4069-881a-16abbe6a420c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:28.2701+00	2025-11-10 08:11:28.2701+00	\N
e52dff32-079c-4d5d-9145-1f70977b29af	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:28.729623+00	2025-11-10 08:11:28.729623+00	\N
f81396c7-a4ef-4c66-bfff-f2f566b7013e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:29.390735+00	2025-11-10 08:11:29.390735+00	\N
29af98b5-433e-49eb-a1a1-bb2e00319fc2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:30.086555+00	2025-11-10 08:11:30.086555+00	\N
68c021c8-a279-457f-9026-bf2766a46bfd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:30.801856+00	2025-11-10 08:11:30.801856+00	\N
1b773c8e-011f-4c90-bf4a-ff35b7adfa99	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:31.570144+00	2025-11-10 08:11:31.570144+00	\N
f0624a57-2e7c-4a66-a19e-aaee5a47694f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:31.936156+00	2025-11-10 08:11:31.936156+00	\N
32550370-ff0b-4e2e-95d3-52cbbe1b44a5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:32.388992+00	2025-11-10 08:11:32.388992+00	\N
3f2b6a2e-797a-4780-ab20-db6c94fba05d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:32.89584+00	2025-11-10 08:11:32.89584+00	\N
e19e3064-f304-4e50-b20b-e413942e79e4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.159	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.128 Mobile/15E148 Safari/604.1	null	2025-11-10 08:11:36.359551+00	2025-11-10 08:11:36.359551+00	\N
f5d6a511-c600-4dd9-b2ed-a337d373b31b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	103.114.174.232	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-10 08:12:36.393196+00	2025-11-10 08:12:36.393196+00	\N
303b84c3-f5f2-45f2-90a2-42fbca8e0c01	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	35ca7248-3189-487d-bf4c-e5157677d040	Created transaction: Tax consultant fee	103.114.174.232	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-10 08:13:03.213188+00	2025-11-10 08:13:03.213188+00	\N
6e8aa5e7-43b7-4d96-99d5-c7abe72aa339	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	d7962291-846d-403e-8843-261e01ff1803	Created transaction: Bike fare	103.114.174.232	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-10 08:13:25.44951+00	2025-11-10 08:13:25.44951+00	\N
658cf521-16dc-45e3-af2e-1a131b9ed748	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	b1f73683-3341-4a4d-bf91-cd7edfbf6f4d	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-10 16:00:41.315351+00	2025-11-10 16:00:41.315351+00	\N
f88a0f22-c23b-438c-90e0-fb2a0fbd8558	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	5bfbe17e-66ce-4513-8cf6-6335bb8c3004	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-10 16:01:27.121032+00	2025-11-10 16:01:27.121032+00	\N
fd95563a-817a-4791-b793-8f9770a1f8dd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-10 17:49:04.433034+00	2025-11-10 17:49:04.433034+00	\N
408951d2-f35a-449f-bd42-bf80b9efe1e2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	2b46ce7c-7c49-4f26-a38f-5a6201c0d81a	Created transaction: Mobile recharge for Maa	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-10 17:50:10.266767+00	2025-11-10 17:50:10.266767+00	\N
ec1b6658-15de-414b-b0ee-781975dede79	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e7573260-58bc-4266-88d6-18bed9968ea2	Created transaction: Water Kettle Repair	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-10 17:50:38.205463+00	2025-11-10 17:50:38.205463+00	\N
fdb98e64-4f14-4219-8702-99f74228b0f8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	2b46ce7c-7c49-4f26-a38f-5a6201c0d81a	Deleted transaction: Mobile recharge for Maa	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-10 17:56:24.176127+00	2025-11-10 17:56:24.176127+00	\N
4daf9b5f-46ee-4f66-a128-68dabdad4706	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	19b78f5a-f034-44aa-b7e4-485ffd8afda5	Created transaction: Singara from Bohera	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-10 17:57:25.318462+00	2025-11-10 17:57:25.318462+00	\N
57ce0705-b904-4375-af02-0df9673d052e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c272af47-ed23-473a-b866-1ef1269d2aaa	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 114, "backfilled": true}	2025-11-09 16:58:14.781816+00	2025-11-09 16:58:14.781816+00	\N
89f4016d-e3f3-4d3f-b829-506e4efb9b50	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	9e4fad59-4d4b-4bf6-97d5-dfc1200ec81f	Created transaction: CNG fare (backfilled)			{"type": "expense", "amount": 200, "backfilled": true}	2025-11-09 16:58:35.862246+00	2025-11-09 16:58:35.862246+00	\N
cbd2e35a-e708-489b-b8d4-b68c4da9dad8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	87b3d61e-80ca-4e14-99fc-1d0763d98efd	Created transaction: Rickshaw fare (backfilled)			{"type": "expense", "amount": 100, "backfilled": true}	2025-11-09 16:58:53.658473+00	2025-11-09 16:58:53.658473+00	\N
7e9f40cf-c6c3-4f47-90bd-98e043e3897c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1394b9c5-a5bb-4d2d-b265-610a264a4c20	Created transaction: TAX pay (backfilled)			{"type": "expense", "amount": 38802, "backfilled": true}	2025-11-09 16:59:45.642014+00	2025-11-09 16:59:45.642014+00	\N
4dcfa066-455e-4112-8ed4-c50abd8acb76	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e75f993c-d932-47f1-b47a-61d116364867	Created transaction: Bike fare (backfilled)			{"type": "expense", "amount": 150, "backfilled": true}	2025-11-09 17:01:16.602348+00	2025-11-09 17:01:16.602348+00	\N
0e5528b8-05c8-4e68-a8a2-23c495a085de	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	af3f64b8-716c-4a1a-85da-d67543d3e588	Created transaction: Bua Bill (backfilled)			{"type": "expense", "amount": 1800, "backfilled": true}	2025-11-09 16:57:42.33479+00	2025-11-09 16:57:42.33479+00	\N
da7196ad-bf8f-4072-a4df-badbd8c8149d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	35ca7248-3189-487d-bf4c-e5157677d040	Created transaction: Tax consultant fee (backfilled)			{"type": "expense", "amount": 2010, "backfilled": true}	2025-11-10 08:13:03.209349+00	2025-11-10 08:13:03.209349+00	\N
9cafc51a-a784-4a33-aa86-d07ef7409b74	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	d7962291-846d-403e-8843-261e01ff1803	Created transaction: Bike fare (backfilled)			{"type": "expense", "amount": 130, "backfilled": true}	2025-11-10 08:13:25.444564+00	2025-11-10 08:13:25.444564+00	\N
ab885bdc-f3c2-49d0-8a27-47e1a153a76e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	b1f73683-3341-4a4d-bf91-cd7edfbf6f4d	Created transaction: Bike fare (backfilled)			{"type": "expense", "amount": 140, "backfilled": true}	2025-11-10 16:00:41.307702+00	2025-11-10 16:00:41.307702+00	\N
8a018271-d82e-480c-8dfb-c405a150f506	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	5bfbe17e-66ce-4513-8cf6-6335bb8c3004	Created transaction: Donation in the path of Allah (backfilled)			{"type": "expense", "amount": 212, "backfilled": true}	2025-11-10 16:01:27.117587+00	2025-11-10 16:01:27.117587+00	\N
05c1785a-0dee-44fa-b96a-d75590565170	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e7573260-58bc-4266-88d6-18bed9968ea2	Created transaction: Water Kettle Repair (backfilled)			{"type": "expense", "amount": 320, "backfilled": true}	2025-11-10 17:50:38.202631+00	2025-11-10 17:50:38.202631+00	\N
dda3e535-3e5c-42b4-a971-572c4405be57	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	19b78f5a-f034-44aa-b7e4-485ffd8afda5	Created transaction: Singara from Bohera (backfilled)			{"type": "expense", "amount": 75, "backfilled": true}	2025-11-10 17:57:25.315287+00	2025-11-10 17:57:25.315287+00	\N
8c4e0999-b8e4-49f4-a9d5-1f11fa727ae7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	be686fc0-a2a1-40b2-bb58-adb7d4e85e65	Created asset: Apple MacBook Pro M1 (backfilled)			{"name": "Apple MacBook Pro M1", "backfilled": true, "purchasePrice": 157000}	2025-11-09 18:04:50.261972+00	2025-11-09 18:04:50.261972+00	\N
08a8164e-a6b3-4f37-b5be-673abdba56b7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	69509677-7d3f-4a00-83a6-509d013171b0	Created asset: Raspberry Pi 4 (backfilled)			{"name": "Raspberry Pi 4", "backfilled": true, "purchasePrice": 11540}	2025-11-09 17:56:39.29489+00	2025-11-09 17:56:39.29489+00	\N
e62918aa-80c2-4114-8ee9-d6ee7aa97f03	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7e82b970-08b0-4ef1-be17-87497b36a457	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-11 14:11:57.291953+00	2025-11-11 14:11:57.291953+00	\N
5ad8a791-4498-4bd2-857e-49514da60ad4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	d9045428-a73b-43a1-9b04-538bfa949944	Created transaction: Bike fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-11 14:12:16.789295+00	2025-11-11 14:12:16.789295+00	\N
734d305c-2c84-4c4c-8488-c536d2359b59	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	b1b86e67-45d7-459c-a136-1c6c7618213c	Created transaction: Bike fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-11 14:12:34.102166+00	2025-11-11 14:12:34.102166+00	\N
c5f870b6-6014-4721-817b-222abb560c85	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	502ddf6d-9935-45b9-881d-28bbd91b8b5c	Created transaction: House rent	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-11 15:55:20.066603+00	2025-11-11 15:55:20.066603+00	\N
9677ebbd-9560-4d45-af0c-3af3cce9cbf3	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-12 15:26:52.541969+00	2025-11-12 15:26:52.541969+00	\N
f26f7cbf-c176-4b24-8584-400e76d0c907	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	176c54b1-3bc7-4d2e-9525-51fa1ffe705c	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-12 15:27:21.812841+00	2025-11-12 15:27:21.812841+00	\N
c7c92ce2-fe62-44ec-bc79-74461259e389	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	dbdb88b3-1c57-4514-b006-a1918b720a29	Created transaction: Snacks	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-12 15:27:36.355483+00	2025-11-12 15:27:36.355483+00	\N
6aed08f8-edbe-4b00-85d1-e050250065eb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	13784bbd-ec16-467e-8d02-a2b5ef9a1a39	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-12 15:27:53.488488+00	2025-11-12 15:27:53.488488+00	\N
87123abb-e814-479a-a403-42bc40d47b45	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	2675b9d5-850a-477e-b734-e48c4a6e2d4a	Created transaction: Rickshaw fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-12 15:28:09.018083+00	2025-11-12 15:28:09.018083+00	\N
5322f7b6-095b-4708-99d0-522750289064	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.64	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.148 Mobile/15E148 Safari/604.1	null	2025-11-13 09:34:33.292604+00	2025-11-13 09:34:33.292604+00	\N
dc78d560-bdce-4bc0-8909-08f31a24c060	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:41:47.452243+00	2025-11-13 14:41:47.452243+00	\N
d301d61e-8137-4392-a634-c2a2017926c7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a2873c77-6b8f-41ac-9cad-4ac3950169cd	Created transaction: Bike fare	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:42:09.6529+00	2025-11-13 14:42:09.6529+00	\N
553a4685-3bf0-4fdb-8627-a281187b1149	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	df345915-6483-4010-8928-32934cde0601	Created transaction: Bohera breakfast	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:42:27.945955+00	2025-11-13 14:42:27.945955+00	\N
33330d7a-40dc-480b-9fb2-c2c7096064c7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	GoalHolding	33283ce5-6dff-4ede-99cd-37335e0870ff	Added holding DPS (Deposit Pension Scheme) to goal bKash DPS - 1783060326999	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:44:10.407201+00	2025-11-13 14:44:10.407201+00	\N
96d66cb0-abb0-4982-b0ee-fdf861f5a1fb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	aa411452-434b-434e-a0ab-ac22ee0c5ea2	Created transaction: Rickshaw fare	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:44:50.844521+00	2025-11-13 14:44:50.844521+00	\N
523b748e-9ef7-41f7-ab45-17acd8eaa5a4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	178abbaf-f077-475f-ac52-7f6bab746ed0	Created transaction: Donation in the path of Allah	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:45:18.779627+00	2025-11-13 14:45:18.779627+00	\N
b5a8cd65-09cd-4b25-af94-56059d7def12	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	0f33a456-55b0-4522-8277-7e77dbcbcddb	Created transaction: Donation in the path of Allah	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:46:17.668762+00	2025-11-13 14:46:17.668762+00	\N
e60cafed-b943-4f48-8d24-9c73ee07bf1e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	176c54b1-3bc7-4d2e-9525-51fa1ffe705c	Updated transaction: Donation in the path of Allah	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:46:34.976181+00	2025-11-13 14:46:34.976181+00	\N
ef06db19-86d7-4b84-bb99-3970d09442c8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	f2ca6c30-e6c3-4222-9213-2eae63c2769c	Created transaction: Snacks	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 14:47:01.755877+00	2025-11-13 14:47:01.755877+00	\N
defd283a-287d-4caa-b3d5-362e5638869e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	508fbd44-d158-43cf-801e-71bbc89f3cc4	Created transaction: Transfer between accounts	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 17:15:22.402286+00	2025-11-13 17:15:22.402286+00	\N
7c1e2b22-4762-445c-a7b4-dff9a13eabc1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	54f40500-5289-4cc6-b4fe-1068debf484a	Created transaction: DESCO Bill Pay	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 17:15:49.627321+00	2025-11-13 17:15:49.627321+00	\N
0f7cbbc2-6e51-4273-a02c-7bc70c05c4ff	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	asset	Asset	928306b1-1bae-42c7-bab7-4c7f9a6b4569	Updated asset: UGREEN 60126 USB-A 2.0 To Type-C 1 Meter Data Cable	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-13 17:16:33.794063+00	2025-11-13 17:16:33.794063+00	\N
8e8749c8-89f3-43e5-b88f-4bbd61fef256	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a959db34-3330-432d-8a8c-14e70b15fbcb	Created transaction: Donation in the path of Allah	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-14 11:28:31.660621+00	2025-11-14 11:28:31.660621+00	\N
828f294f-3f0b-4303-aa26-10443b1cd11b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	asset	Asset	e9e06ac5-2f52-4db5-8227-3ef3201b56c9	Created asset: iPhone 13 Mini	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-14 11:29:34.617898+00	2025-11-14 11:29:34.617898+00	\N
785126a2-9a19-47ab-8f61-9f201d20a87c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	asset	Asset	e9e06ac5-2f52-4db5-8227-3ef3201b56c9	Updated asset: iPhone 13 Mini	202.181.7.21	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-14 11:29:54.054836+00	2025-11-14 11:29:54.054836+00	\N
a799e151-971b-4fa1-8a53-c63f1bb5faa7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-14 16:02:37.257632+00	2025-11-14 16:02:37.257632+00	\N
dccf4eae-10fa-4341-89e6-80bc5ef59c9f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	d9420559-1b1b-4907-9c55-7b1acbc42833	Created transaction: Transfer from Rocket to bKash via NPSB	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-14 16:03:12.684966+00	2025-11-14 16:03:12.684966+00	\N
f13b13bc-5c62-4ea9-ac66-894a70bc7c74	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	RecurringTransaction	8845435a-f95f-448b-850a-dfb53371cf08	Updated recurring transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 08:02:33.184892+00	2025-11-15 08:02:33.184892+00	\N
7ff5891d-ab0c-4ebb-8611-a55def4f9195	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	2fe39b89-056e-4f53-bfb0-0e63d65ceef9	Deleted transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 08:03:00.538911+00	2025-11-15 08:03:00.538911+00	\N
2abc3ae3-f0ce-46f5-bdda-dfba09dbc8e4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	RecurringTransaction	8845435a-f95f-448b-850a-dfb53371cf08	Deleted recurring transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 08:03:08.151659+00	2025-11-15 08:03:08.151659+00	\N
29fcd3c1-07cc-49cf-a125-b6ceba94e38f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	RecurringTransaction	8845435a-f95f-448b-850a-dfb53371cf08	Deleted recurring transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 08:03:14.66768+00	2025-11-15 08:03:14.66768+00	\N
26548b9f-23cf-4180-94b1-e6666db4a9c5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	ef7e519d-2cb1-4e6f-9750-2ac80dacd4dc	Created transaction: Tehari Ghor Dhanmondi	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:11:26.487188+00	2025-11-15 15:11:26.487188+00	\N
490f4808-b80f-45dc-b041-7da33a60a4ca	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	ca5c220a-4d8e-4618-a639-4e11c1202aa4	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:11:58.447117+00	2025-11-15 15:11:58.447117+00	\N
95f22f76-1a80-42ca-a260-ada2484b2e19	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	bc2d45c1-fe27-48e8-abae-5b53c9b929bf	Created transaction: Bike fare	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:12:23.421583+00	2025-11-15 15:12:23.421583+00	\N
e1c6c3d3-1e62-421d-a50a-3eee5e4b2dd1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e6e57da4-2aa6-4b5b-9e7f-04c1bf0180a2	Created transaction: CNG fare	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:13:14.600886+00	2025-11-15 15:13:14.600886+00	\N
7961d7f8-78a9-495b-8464-5541d2e285b8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7dfa942c-2897-42d8-83d9-ee3ddb69bac4	Created transaction: From foodpanda	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:15:03.008803+00	2025-11-15 15:15:03.008803+00	\N
8c43bc53-eaa5-44ad-9a47-8935c587dbea	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	372607e8-c2ba-4620-a9b7-6e0fecb035a8	Created transaction: From foodpanda	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:20:30.440635+00	2025-11-15 15:20:30.440635+00	\N
14615e33-e354-4c5e-9074-d4df076c66c5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	4d63308a-ea1c-465e-9a35-eb5d8d47f60f	Created transaction: From Daraz	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:22:58.599412+00	2025-11-15 15:22:58.599412+00	\N
325375e6-1db9-4a66-b247-0cceba67f564	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3f6106fa-36fc-491e-91c5-4a9248f3a334	Created transaction: From aarong	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:23:17.639734+00	2025-11-15 15:23:17.639734+00	\N
b79590b7-0dca-4c42-95d6-21bf089244ae	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	a27ef7d5-3878-42ab-a1ba-e004363051e7	Deleted transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:30:02.761484+00	2025-11-15 15:30:02.761484+00	\N
76c2f5e2-d1e0-4650-b244-85c5b806d1c4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	176c54b1-3bc7-4d2e-9525-51fa1ffe705c	Deleted transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:30:09.226685+00	2025-11-15 15:30:09.226685+00	\N
4bfa07d8-0331-40c5-a2b2-259570eb7e54	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	7e82b970-08b0-4ef1-be17-87497b36a457	Deleted transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:30:13.747363+00	2025-11-15 15:30:13.747363+00	\N
396527e5-83cb-47b6-810d-5091783d03ef	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	e625c60a-636d-4573-9578-dbe2e6315320	Deleted transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:30:29.781086+00	2025-11-15 15:30:29.781086+00	\N
9f6b1227-113c-4fc7-bed0-949421591d5e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	3472a26a-046a-45e9-823e-007907a41983	Deleted transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:30:52.208384+00	2025-11-15 15:30:52.208384+00	\N
8c38f4c2-dac8-4e8a-ab6b-8bb5b71fbea9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	ce9b0658-afba-4435-90be-318cece06e26	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-16 17:09:55.812537+00	2025-11-16 17:09:55.812537+00	\N
3f01415b-aa8d-479b-82ec-78164518bab2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1e78f3b2-8f2e-474d-941f-4a3573a0e3f3	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-16 17:10:12.54099+00	2025-11-16 17:10:12.54099+00	\N
9a995c27-0ec0-49d5-924e-ae0e13a736cb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e5b30afb-5004-48c1-bd52-7318a3dbea5f	Created transaction: Snacks	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-16 17:10:26.86669+00	2025-11-16 17:10:26.86669+00	\N
40489682-014e-44a1-aaa4-dd12f992db95	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a0464f2d-ac2d-4d6e-b449-8c8d66e66a81	Created transaction: Bike fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-16 17:10:38.882134+00	2025-11-16 17:10:38.882134+00	\N
a1e4090e-48e5-418d-84d3-6241306b0894	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e47dca06-1702-4b5b-9983-0204f10f69b7	Created transaction: From Vivasoft for KPI	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-16 17:11:00.89823+00	2025-11-16 17:11:00.89823+00	\N
d5d48b8a-049d-44f9-9ff1-e940bcda3750	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	4f614975-debf-4fbf-b1d7-b273280c7f79	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-16 17:11:35.563866+00	2025-11-16 17:11:35.563866+00	\N
d15f99d3-7863-4ad9-aa1a-570edd98beb5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	delete	transaction	Transaction	c6245739-4004-4465-b48f-eb0752d2a93b	Deleted transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-15 15:31:55.466774+00	2025-11-15 15:31:55.466774+00	\N
7153e0e1-487d-4796-b421-03ca91a9ad2e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-16 17:07:00.076989+00	2025-11-16 17:07:00.076989+00	\N
7b2173ce-15c4-48e0-8c51-14fbdd6c91aa	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	5c60f9b8-7a6d-4a17-9a0a-f5b8613b3b1c	Created transaction: Bike fare	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 15:00:29.708408+00	2025-11-17 15:00:29.708408+00	\N
a7ce2009-ba70-486d-afc5-1b7ad2b1aae0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	f3778271-821f-43bd-9563-acf5cbc01ab0	Created transaction: Rickshaw fare	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 15:00:46.360357+00	2025-11-17 15:00:46.360357+00	\N
4d66bf2c-cafd-4ad8-b548-aa38ddc61218	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	5c60f9b8-7a6d-4a17-9a0a-f5b8613b3b1c	Updated transaction: Bike fare	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 15:00:52.8966+00	2025-11-17 15:00:52.8966+00	\N
35133b26-d3e1-4964-99bb-a719449f851b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	23ea0999-f751-4b4c-bd42-5ac0f2ffbdac	Created transaction: Gift for humaun bhai	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 15:01:38.482012+00	2025-11-17 15:01:38.482012+00	\N
ce7d1ca5-88cb-481b-b33b-25dd09274b00	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7fbfc35d-0a29-4f65-bc2b-9aa613c5da05	Created transaction: bKash Staff for go to ST	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 15:02:03.021906+00	2025-11-17 15:02:03.021906+00	\N
ac6d24dd-d816-4dbc-aef8-0e775767811c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	26940eb4-ae0b-44b1-9834-2fb7954da753	Created transaction: bKash HR for night activity	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 15:02:31.074722+00	2025-11-17 15:02:31.074722+00	\N
e362afce-70fd-4b71-8fba-36077ec591d3	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	07271880-ccf5-402c-b31f-73337e4553cc	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 15:03:00.410178+00	2025-11-17 15:03:00.410178+00	\N
8a7e200e-b7dd-4e40-abef-42b6a73eb8b2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-17 17:37:31.481292+00	2025-11-17 17:37:31.481292+00	\N
aabbb11e-5622-42cd-bec5-83e74d3855c5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e7f9fd22-5639-448a-9582-47c1132334eb	Created transaction: Snacks	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 15:48:19.905039+00	2025-11-18 15:48:19.905039+00	\N
3fe68d78-2a93-413b-84b7-f1fd21cde515	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	6e521053-d2cd-41c8-8b3c-a3ee93188d3b	Created transaction: Rickshaw fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 15:48:33.905629+00	2025-11-18 15:48:33.905629+00	\N
7955a8e1-9706-4046-9311-66f73bac9668	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c80eface-97e5-4f39-aedc-c5ac489f148c	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 15:48:52.107583+00	2025-11-18 15:48:52.107583+00	\N
8565eac0-74bb-4fe8-bff2-ebe56462936d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	6e521053-d2cd-41c8-8b3c-a3ee93188d3b	Updated transaction: Rickshaw fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 15:48:59.508407+00	2025-11-18 15:48:59.508407+00	\N
c990172e-301c-44c5-a4ac-e7a77b6b2b1e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	d3860014-cd06-486d-98fc-96bd85a5b809	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 15:49:30.472852+00	2025-11-18 15:49:30.472852+00	\N
e1c3fed1-c8f1-4dae-b14b-0b5d2be790c0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	584ce5de-1559-4e95-94a9-1d2d7a6ea7a2	Created transaction: Credit card yearly fee	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 15:50:00.231485+00	2025-11-18 15:50:00.231485+00	\N
ca46f2fd-9968-4b97-9a1a-7c0916115a07	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	goal	GoalHolding	28655439-8c56-4ed1-82b6-b0db265301e1	Added holding DPS (Deposit Pension Scheme) to goal bKash DPS - 1783060302607	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 16:03:01.06189+00	2025-11-18 16:03:01.06189+00	\N
48624fe7-7d76-4b0a-90ad-a396d0a13816	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	fa02af77-895d-480a-8993-2fa2e59c71c6	Created transaction: Transfer between accounts	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 17:04:51.686957+00	2025-11-18 17:04:51.686957+00	\N
8fe0f764-4630-490d-97aa-e9a2a864fce6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	72d54b38-9878-414d-967e-14eb55d1578d	Created transaction: Cashback	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 17:05:14.339777+00	2025-11-18 17:05:14.339777+00	\N
28bebb54-8f91-4400-8fb1-46e7ec29babc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7afa2bca-5952-4ad4-99b0-056a7e620803	Created transaction: Medicine from Aroggo	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-18 17:05:43.06204+00	2025-11-18 17:05:43.06204+00	\N
4d7bbacc-4ac2-412b-8dc8-213de920d67b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	103.114.174.232	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0	null	2025-11-19 11:57:17.924391+00	2025-11-19 11:57:17.924391+00	\N
0cd5df50-52e4-458c-9d2c-873e2669f18e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.87	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.148 Mobile/15E148 Safari/604.1	null	2025-11-20 11:34:09.809252+00	2025-11-20 11:34:09.809252+00	\N
417d23e5-0e11-42cf-8e10-ab6263d767fd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.123	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/142.0.7444.148 Mobile/15E148 Safari/604.1	null	2025-11-20 13:15:41.554167+00	2025-11-20 13:15:41.554167+00	\N
6266be84-be9e-4d40-8886-d27225c2c5e6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:54:48.875813+00	2025-11-20 14:54:48.875813+00	\N
b1e79935-c09f-49de-87ad-94428ed15de7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	f89aa192-747d-4cba-8176-8b3758618a46	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:55:19.686708+00	2025-11-20 14:55:19.686708+00	\N
7da180a8-39e1-45be-88ca-56da1302e8b6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	804ed2b9-d8ae-4d2c-b1ef-6e02a0909878	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:55:40.738774+00	2025-11-20 14:55:40.738774+00	\N
3ed148dd-eb43-42c8-9461-5e149883c68b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1901245e-e128-4cca-9f89-034b42a94946	Created transaction: Bike fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:55:53.390466+00	2025-11-20 14:55:53.390466+00	\N
ee7aeec0-cf15-4394-a30a-6ae12c72f196	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	ec2a230a-ca06-49af-974a-529a62470d57	Created transaction: Rickshaw fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:56:18.259091+00	2025-11-20 14:56:18.259091+00	\N
a30528b5-eb59-4758-858c-1c7b82ce98b9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	477b51d9-636f-4b7c-9ac0-605d0a40d81e	Created transaction: Rickshaw fare	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:56:35.560221+00	2025-11-20 14:56:35.560221+00	\N
571030d2-10fe-4302-8cc5-69f75ae781ee	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3b255be4-a638-4a9c-bcb0-de1c8b41ff77	Created transaction: Snacks	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:56:53.561592+00	2025-11-20 14:56:53.561592+00	\N
ac1c6e23-4bb9-4f2d-8124-b1c5e32be71a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	f4b9815e-51ad-4ff8-99ad-3900c05626a8	Created transaction: Snacks	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:57:07.869721+00	2025-11-20 14:57:07.869721+00	\N
d8a99edd-d91a-4972-855e-b0866427a943	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3c4eee26-56d3-42ce-8a41-338d8357f42f	Created transaction: Bohera	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:58:27.958003+00	2025-11-20 14:58:27.958003+00	\N
40fc7a93-8266-471f-8ff0-889e2bb3bf1b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	fee68dd4-f949-4548-a6dd-816e343096af	Created transaction: Transfer between accounts	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:58:40.402572+00	2025-11-20 14:58:40.402572+00	\N
c78d6cde-8cc1-4175-8817-f1b50f6a2703	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a9b1999d-43f1-4ea1-8c23-d28aee9f1bed	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-20 14:59:21.263609+00	2025-11-20 14:59:21.263609+00	\N
b4609c93-8393-46b3-85ba-ff43293f8213	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:48:07.501558+00	2025-11-22 16:48:07.501558+00	\N
9ac4fb49-2451-46f6-b0fb-7dde0e37f983	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	622d1295-9608-49e1-ba8c-2d3ad25605c5	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:49:33.212584+00	2025-11-22 16:49:33.212584+00	\N
ae302f3b-6e44-4774-9805-37eba2c600b5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	dbf29958-e708-445b-941b-b855c975b2f6	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:49:46.36528+00	2025-11-22 16:49:46.36528+00	\N
df7d37a1-7e5b-4be6-9047-51ce45b79466	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a60047ea-9997-4169-bb55-0cef6e62f636	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:49:57.755649+00	2025-11-22 16:49:57.755649+00	\N
2782a7fd-af65-4861-bb0b-9ba9acc765f5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	90a92532-7788-42c8-a978-08783d158143	Created transaction: Transfer between accounts	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:50:19.414952+00	2025-11-22 16:50:19.414952+00	\N
cc96924f-37af-4e54-ad9f-c13def800fd4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	lend	Lend	9bf5b5ac-6514-41e0-bdb3-5328c6970dd1	Created lend record for Nurul Huda	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:51:06.156414+00	2025-11-22 16:51:06.156414+00	\N
98549e6d-b77c-4e6a-8d4e-e70a1d51c84f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	397fab37-18fd-4982-8b10-ed97a3956141	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:51:33.226228+00	2025-11-22 16:51:33.226228+00	\N
71f94bf5-eea2-4051-9821-fbc67458f252	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	8ab6b900-bce6-4604-a9f3-4ff4d4eb74e9	Created transaction: Snacks	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:51:49.231857+00	2025-11-22 16:51:49.231857+00	\N
0e5766f3-1595-430d-8442-35310b3afd19	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	bb332524-84ff-4995-82fd-5abd92866552	Created transaction: Bus fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-22 16:52:04.704618+00	2025-11-22 16:52:04.704618+00	\N
e005385a-bc04-4b3b-9f36-68cb802d3af9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	851dc38a-6c3d-4f55-96ae-b0301a9cc6dc	Created transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 16:19:54.641309+00	2025-11-23 16:19:54.641309+00	\N
6ba4db3f-3e35-42f4-8536-e048d722296b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	fa731082-60f7-4fa3-9615-e7cf107f9e01	Created transaction: Snacks	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 16:20:33.320033+00	2025-11-23 16:20:33.320033+00	\N
0d67e789-8be3-4a53-b4da-978571127d94	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	483a1bb7-aba0-4486-b176-5aaba95b5327	Created transaction: Imsurence	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 16:22:05.086347+00	2025-11-23 16:22:05.086347+00	\N
8ba71368-091e-4349-a928-1bb24a0cc26a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	cef0a836-ff4f-46ef-acfe-a43212d253ba	Created transaction: Insurence	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 16:22:43.18813+00	2025-11-23 16:22:43.18813+00	\N
3fe2fa88-70d5-46d6-8628-45621cc647d4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	aa1bf759-659b-4379-b7fc-190b5748e083	Created transaction: Transfer between accounts	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 16:27:32.973694+00	2025-11-23 16:27:32.973694+00	\N
d11b4d0f-e0be-4dd5-838a-496f1c857332	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 17:07:31.067133+00	2025-11-23 17:07:31.067133+00	\N
0c8a2497-e627-4bd8-9dad-c1a4119cd84a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	cef0a836-ff4f-46ef-acfe-a43212d253ba	Updated transaction: Insurance	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 17:32:48.828451+00	2025-11-23 17:32:48.828451+00	\N
7e93b6f3-7bd7-4b3c-a20c-e1473d66fc3c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	f4b9815e-51ad-4ff8-99ad-3900c05626a8	Updated transaction: Snacks	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 18:10:31.991234+00	2025-11-23 18:10:31.991234+00	\N
14f03581-d9c3-4799-bcb4-db11e2328ed8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	cef0a836-ff4f-46ef-acfe-a43212d253ba	Updated transaction: Insurance	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-23 18:20:06.871942+00	2025-11-23 18:20:06.871942+00	\N
1fdfd6cf-1526-49ab-9e8c-45cf2ec60bab	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	37.111.194.241	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/143.0.7499.38 Mobile/15E148 Safari/604.1	null	2025-11-24 03:41:22.387397+00	2025-11-24 03:41:22.387397+00	\N
64a0380b-f652-41df-aa26-f9b806ea8d6b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7c505ec8-bdb2-4898-a6fe-7d1d8f30b3ff	Created transaction: Rickshaw fare	37.111.194.241	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/143.0.7499.38 Mobile/15E148 Safari/604.1	null	2025-11-24 03:41:55.130218+00	2025-11-24 03:41:55.130218+00	\N
0a8a781a-dff5-4d2e-93cd-ae26d732d40a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	eaaf61a5-c5fe-47f8-a2f3-260d2b43a151	Created transaction: Rickshaw fare	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 15:00:53.05951+00	2025-11-24 15:00:53.05951+00	\N
2d884407-6f3c-4fa6-940c-9e9e7b66c02b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	4e618279-5e38-47b3-91a2-e2f111270099	Created transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 15:01:12.497935+00	2025-11-24 15:01:12.497935+00	\N
9868f079-2a65-4919-94c8-6c65b32ebfc0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	02a535dd-be27-41dd-8bc3-4991f9886ac1	Created transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 15:01:24.989124+00	2025-11-24 15:01:24.989124+00	\N
3dbd5257-fb5e-4755-853c-45743988a7cd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	dfb32f1c-9d8c-4376-9f54-27362e786b9a	Created transaction: Donation in the path of Allah	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 15:16:03.076129+00	2025-11-24 15:16:03.076129+00	\N
f9f4ccfd-1df6-4d90-b95b-2e5a639d0cfd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	12adafa0-f4b9-4ad3-85f0-83daa03477dc	Created transaction: Mess bajar	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 15:57:15.022458+00	2025-11-24 15:57:15.022458+00	\N
2878af65-978c-4686-8a40-3c201d89dcc2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a56ce51d-07a7-4f8b-b109-94f995380011	Created transaction: Mess Bazer	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 15:57:48.253239+00	2025-11-24 15:57:48.253239+00	\N
4a6f0734-e5b0-44cb-bc16-a85bd7f7c046	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	81bb1c6f-4850-45ba-a793-0f9a0f8859bc	Created transaction: Fruits	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 15:58:07.805267+00	2025-11-24 15:58:07.805267+00	\N
4fe41848-c87a-4a05-93ef-02ee61c6ca5d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	692eee26-d26d-4dd9-9207-70982b6ac8ea	Created transaction: Mess Bazar	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 16:01:51.478578+00	2025-11-24 16:01:51.478578+00	\N
06ce0317-f1b0-4fca-935c-1d8c1aad2d40	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	36372654-9013-40dd-960c-ee9eb5c84f58	Created transaction: Fruits	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 16:02:10.456115+00	2025-11-24 16:02:10.456115+00	\N
5bc19e5b-60ae-43a3-81ff-c63a030e28c7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	credit_card	CreditCardPayment	d54a92d0-35c4-47fa-abb4-c3dc6de3a344	Recorded credit card payment: Brac Bank (Credit Card)	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-24 16:03:22.256493+00	2025-11-24 16:03:22.256493+00	\N
abadb2ca-5536-4d94-844b-481c28d8a04c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-25 15:51:54.00882+00	2025-11-25 15:51:54.00882+00	\N
2ee5a627-526a-43ee-a94e-fdb4ebf292d1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	285d505a-fb8c-42da-ac0a-4cf79b06af3c	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-25 15:52:34.359195+00	2025-11-25 15:52:34.359195+00	\N
a5e84201-3166-490b-ab78-f5d5d4dcfe82	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c9f7943a-69ee-4461-86e5-83fe8f5b18c5	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-25 15:52:47.200346+00	2025-11-25 15:52:47.200346+00	\N
4a20eed4-5e89-48ed-a2ec-a98d244bf100	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	42024371-e47e-4464-ae20-4b759566c5c0	Created transaction: Donation in the path of Allah	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-25 15:55:35.697066+00	2025-11-25 15:55:35.697066+00	\N
0bd9272e-cb07-4ce4-879e-8ebf70154b30	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7874feb8-44eb-42e5-a3f4-1ee9f99630ce	Created transaction: Nov 2025/Salary	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-25 15:56:13.766976+00	2025-11-25 15:56:13.766976+00	\N
5eae23ed-6ca8-400a-859f-19289703d470	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	7d695631-2ebb-4aec-985d-77078cb350cf	Created transaction: Nov 2025/Salary	104.28.240.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-25 15:56:30.595547+00	2025-11-25 15:56:30.595547+00	\N
ef416f76-6c4c-416f-8685-9b145c8839c3	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:29:23.359549+00	2025-11-27 16:29:23.359549+00	\N
cff7d49b-9478-4e72-91a7-d804183adfe8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	31a15c49-17e9-4e17-8036-77f7e6ef27a2	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:29:54.976947+00	2025-11-27 16:29:54.976947+00	\N
42404e29-291e-4c8e-af8c-4edf48a3d5d6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	4ca91ce0-6e8c-43de-b611-d3caaf07fe84	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:30:08.766124+00	2025-11-27 16:30:08.766124+00	\N
4af74871-e98a-4411-ac7a-f858b8f31d7c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	33cd2521-0b2b-4556-ad23-a7701ff7d747	Created transaction: Rickshaw fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:30:18.661867+00	2025-11-27 16:30:18.661867+00	\N
3f330931-b44b-4aa5-ac3a-a64e6e5875b2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	f5948d5c-b2c7-46b2-a12e-78b18f635988	Created transaction: Bike fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:30:31.386861+00	2025-11-27 16:30:31.386861+00	\N
e26671a1-4069-4b31-87fb-8ddec6183551	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	57ad2986-8eb1-4ced-a93d-2d38b2b07706	Created transaction: CNG fare	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:30:50.149735+00	2025-11-27 16:30:50.149735+00	\N
c51b9db6-bfee-4877-895d-f5436791715f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	4964d3cf-b729-493e-b039-ede372efd9e1	Created transaction: Insurance	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:31:47.314793+00	2025-11-27 16:31:47.314793+00	\N
3447106e-2733-49d0-ac08-e09f886eaed6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	f5088e9d-ce7e-496c-af7a-e9245b31caba	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:32:11.539259+00	2025-11-27 16:32:11.539259+00	\N
cfac469a-a60f-444e-93b5-967d9bf18c7f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	5456aa2f-1a77-4fd6-a289-7b0f7591a822	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:32:21.910063+00	2025-11-27 16:32:21.910063+00	\N
c7904c30-2ebe-46cc-aaf5-e524a961795b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	c45fd1ab-8363-463a-ae36-bd41d8af0406	Created transaction: Donation in the path of Allah	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:32:31.734844+00	2025-11-27 16:32:31.734844+00	\N
d4b00f28-7d6f-45d2-aab0-485402ead781	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3a43a967-62f2-4797-8b09-f71579a96537	Created transaction: Mess Bazar	104.28.208.84	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-27 16:33:23.082216+00	2025-11-27 16:33:23.082216+00	\N
81f208f9-6f53-423f-8f27-46b52b0b00fd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:09:15.535858+00	2025-11-29 03:09:15.535858+00	\N
fdb3d695-76c1-4992-9b68-b92b237bb4ff	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	b8d126d1-b9cc-4ef5-97fa-e8fdb9005476	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:10:01.431696+00	2025-11-29 03:10:01.431696+00	\N
4b3b7086-4cb9-4986-a609-764897b0b234	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	52b975d5-1a20-4ef4-9c22-7b1350214e6a	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:10:26.344705+00	2025-11-29 03:10:26.344705+00	\N
31a75faf-6468-4e59-b2f6-1735e6aed0fa	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	e45eb6d3-e49d-4ee4-8b07-fd8f49fc4d63	Created transaction: Bike fare	104.28.208.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:10:53.978315+00	2025-11-29 03:10:53.978315+00	\N
91681374-c1fb-4eeb-9521-c4e62bd1a3b1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	8f57d6b9-5eec-404e-a4a9-e8737d82363e	Created transaction: Donation in the path of Allah	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:16:31.373688+00	2025-11-29 03:16:31.373688+00	\N
d5444ead-3bd9-4fc1-b831-4f95b29d38ca	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	1be42af4-7828-4338-aaf3-3292b87f1d9e	Created transaction: Donation in the path of Allah	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:17:02.342721+00	2025-11-29 03:17:02.342721+00	\N
9ce22185-ce3b-4c94-86fe-52c92075e293	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	83464a1c-f014-4bfc-99b2-1ae5a7c823c4	Created transaction: Donation in the path of Allah	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:18:19.065798+00	2025-11-29 03:18:19.065798+00	\N
c5550356-4373-475d-b4bb-e29df0a86d58	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	a22ea27c-373f-4a19-b32d-483d75005f82	Created transaction: Mobile recharge	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:18:50.050523+00	2025-11-29 03:18:50.050523+00	\N
21944579-851b-40d5-86f5-ee1408049bd7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	8740f4f9-b89d-40dd-9ae0-144ac0a6a4e0	Created transaction: Medicine	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:19:26.980989+00	2025-11-29 03:19:26.980989+00	\N
e7b87ad2-cea5-41ac-8d5f-621b885a6af5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	56b94022-7c6f-40bc-87e2-d88b54e1ded0	Created transaction: Donation in the path of Allah	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:20:03.482235+00	2025-11-29 03:20:03.482235+00	\N
1be32a09-c957-47ca-a979-8f0f20149680	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	38b3e6d9-fc0e-4026-9b99-f55d3f33a163	Created transaction: Snacks	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:20:36.500887+00	2025-11-29 03:20:36.500887+00	\N
3ca7b434-dac8-413a-a46e-d398f79e6dc8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	479aafc0-7179-4066-9273-ca129f956e26	Created transaction: Raqi	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:21:40.25047+00	2025-11-29 03:21:40.25047+00	\N
5bdb0e24-17ca-4b37-a73f-5eb7845f455b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	0fbf67b0-0c14-48a4-adfe-0576c0193167	Created transaction: From Foodpanda	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:22:48.035688+00	2025-11-29 03:22:48.035688+00	\N
66cf99f7-d25b-46b8-b524-156b8a830378	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3c763210-e1f2-4ee6-853a-4bb93f0ae105	Created transaction: For mess	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:23:30.287896+00	2025-11-29 03:23:30.287896+00	\N
cd83a4c5-ed45-45b4-ae42-aeef0d452a71	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	33c98a2c-b880-4b85-beb4-e6d453283839	Created transaction: For mess	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:24:56.763694+00	2025-11-29 03:24:56.763694+00	\N
393d768f-5e7f-403b-8455-ee618efd1c9d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	lend	Lend	a19e8944-a810-4272-a46e-2fe47fdcb35f	Created lend record for Jamal bhai	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:25:54.120036+00	2025-11-29 03:25:54.120036+00	\N
0ea08d28-eed8-4e23-ba0f-07dce4436e63	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	3a60ec7f-30ff-452a-aa16-5c3b96c0ebf2	Created transaction: From Lukman bhai for Youtube	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:30:36.316971+00	2025-11-29 03:30:36.316971+00	\N
f71de3e5-cfa0-4cd6-a7a9-5e5ad17578ab	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	update	transaction	Transaction	b28574d2-c2ab-4cae-9b04-cd9007826c4b	Updated transaction: Lent to Jamal bhai	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:30:50.471191+00	2025-11-29 03:30:50.471191+00	\N
4f453177-29bd-4137-a043-ac834403661d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	create	transaction	Transaction	72e1f60f-fb50-4ff8-8f3d-b01b6e058bd6	Created transaction: Cashback	104.28.240.85	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36	null	2025-11-29 03:37:05.41221+00	2025-11-29 03:37:05.41221+00	\N
3360952f-df5c-43b4-aca5-0ee3de2ba39f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	202.181.7.21	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/143.0.7499.38 Mobile/15E148 Safari/604.1	null	2025-11-29 05:28:52.095094+00	2025-11-29 05:28:52.095094+00	\N
8089721f-a6bd-49cf-9c40-2f1622c9a2f8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	login	auth	User	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	User logged in successfully	202.181.7.21	Mozilla/5.0 (iPhone; CPU iPhone OS 26_1_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/143.0.7499.38 Mobile/15E148 Safari/604.1	null	2025-11-29 05:30:58.907363+00	2025-11-29 05:30:58.907363+00	\N
\.


--
-- Data for Name: asset_attachments; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.asset_attachments (id, user_id, asset_id, file_name, original_name, file_path, file_url, file_size, mime_type, attachment_type, description, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: assets; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.assets (id, user_id, name, description, category, brand, model, serial_number, purchase_date, purchase_price, purchase_location, warranty_start_date, warranty_end_date, warranty_provider, warranty_type, status, notes, created_at, updated_at, deleted_at) FROM stdin;
a39f51d8-52a8-4fd5-b2dd-6e01f27fe0a9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	GRID Monitor Stand 		Furniture	GRID	White		2023-05-12 00:00:00	2500	https://99grid.com	\N	\N			active		2025-11-08 14:46:35.64897+00	2025-11-08 14:46:35.64897+00	\N
3773b103-2b4e-4a97-92bc-02711c7dfe3d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	GRID Industrious Work Desk		Furniture	GRID	Light Wood Grain / Compact: 36”W x 24”D x 30”H	Order #12330	2023-04-29 00:00:00	5000	https://99grid.com	\N	\N			active		2025-11-08 15:19:09.007209+00	2025-11-08 15:19:09.007209+00	\N
8fce39d8-3365-429c-a95e-f4c220d92e30	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	GRID Filing Cabinet	Order #12330	Furniture	GRID	Light Wood Grain / Wheel	Order #12330	2023-04-29 00:00:00	5500	https://99grid.com	\N	\N			active		2025-11-08 15:21:39.477664+00	2025-11-08 15:21:39.477664+00	\N
538e0600-4050-45c7-bd63-0193baabaf47	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	GRID Comfy Chair		Furniture	GRID	Black	Order #15938	2023-11-18 00:00:00	16500	https://99grid.com	\N	\N			active		2025-11-08 15:23:29.93397+00	2025-11-08 15:23:29.93397+00	\N
c701367f-53d9-4354-96a3-dc982e0eca9a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Seagate STHH2000400 Backup plus Ultra touch 2TB USB-C Portable Hard Drive		Electronics	Seagate	Ultra touch 2TB USB-C	STHH2000400	2023-09-21 00:00:00	9000	https://www.startech.com.bd	\N	\N			active		2025-11-08 15:33:25.891351+00	2025-11-08 15:33:25.891351+00	\N
ed6a60bc-86f8-426d-af91-466e70edf585	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Ugreen PB760 10000mAh 20W Magnetic Wireless Power Bank #35341	Order##666245	Electronics	Ugreen	10000mAh 20W Magnetic Wireless Power Bank	PB760	2025-06-05 00:00:00	4020	https://www.startech.com.bd	\N	\N			active		2025-11-08 16:00:05.718367+00	2025-11-08 16:00:05.718367+00	\N
c8256dd6-c1fa-4fd0-959b-64d3fec1d2fe	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Havit H94 4-Port High-Speed USB Hub	Order#437463	Electronics	Havit	4-Port	H94	2024-03-29 00:00:00	1025	https://www.startech.com.bd	\N	\N			active		2025-11-08 16:02:10.069017+00	2025-11-08 16:02:10.069017+00	\N
6ce7eda6-068a-4074-a33e-2047d589d9fe	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Transcend 64GB Micro SD UHS-I U1 Memory Card with Adapter (TS64GUSD300S-A)	Order#435920		Transcend	64GB	TS64GUSD300S-A	2024-03-24 00:00:00	800	https://www.startech.com.bd	\N	\N			active		2025-11-08 16:05:19.589561+00	2025-11-08 16:05:19.589561+00	\N
e78b4635-1c7d-4f38-b76d-947ae97ee11d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	UiiSii HM13 Wired In-Ear Headphone with Mic	Order#94494		UiiSii		HM13	2021-01-20 00:00:00	730	https://www.startech.com.bd	\N	\N			active		2025-11-08 16:14:36.62371+00	2025-11-08 16:14:36.62371+00	\N
7725934c-624b-413d-adac-40bde09e1cc1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	TP-Link Deco M4 (Single Pack) Whole Home Mesh Wi-Fi System AC1200 Dual-band Router	Order#247961		TP-Link	Deco M4	AC1200	2022-10-02 00:00:00	5200	https://www.startech.com.bd	\N	\N			active		2025-11-08 16:33:04.542262+00	2025-11-08 16:33:04.542262+00	\N
1ffb3514-d3b2-4674-9033-cfe29856d983	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Apple AirPods Pro	Order#0304364565		Apple	Pro		2022-04-04 00:00:00	21990	https://www.pickaboo.com	\N	\N			active		2025-11-08 16:38:31.90546+00	2025-11-08 16:38:31.90546+00	\N
be686fc0-a2a1-40b2-bb58-adb7d4e85e65	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Apple MacBook Pro M1		Electronics	Apple	MacBook Pro M1		2021-04-14 00:00:00	157000	https://apple.com	\N	\N			active		2025-11-09 18:04:50.261972+00	2025-11-09 18:04:50.261972+00	\N
69509677-7d3f-4a00-83a6-509d013171b0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Raspberry Pi 4	Order#178218	Electronics	Raspberry	Pi 4		2020-10-31 00:00:00	11540	https://techshopbd.com	\N	\N			active		2025-11-09 17:56:39.29489+00	2025-11-09 18:09:45.422277+00	\N
928306b1-1bae-42c7-bab7-4c7f9a6b4569	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	UGREEN 60126 USB-A 2.0 To Type-C 1 Meter Data Cable	Order#437463	Electronics	UGREEN	USB-A 2.0 To Type-C 1 Meter 	60126	2024-03-29 00:00:00	325	https://www.startech.com.bd/	\N	\N			active		2025-11-08 16:03:18.026248+00	2025-11-13 17:16:33.790932+00	\N
e9e06ac5-2f52-4db5-8227-3ef3201b56c9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	iPhone 13 Mini		Electronics	Apple	13 Mini		2022-03-27 00:00:00	63000	Medina	\N	\N			active		2025-11-14 11:29:34.612983+00	2025-11-14 11:29:54.051276+00	\N
\.


--
-- Data for Name: bill_payments; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.bill_payments (id, user_id, bill_id, amount, payment_date, account_id, notes, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: bills; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.bills (id, user_id, name, category, amount, frequency, start_date, due_day, last_paid_date, last_paid_amount, auto_pay, reminder_days, active, notes, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: budgets; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.budgets (id, user_id, category_id, amount, period, custom_start_date, custom_end_date, rollover, alert_threshold, enabled, notes, created_at, updated_at, deleted_at) FROM stdin;
f6f03cd1-a6a0-4406-8a0b-dc8f98935c15	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	transport	5000	monthly	\N	\N	f	80	t		2025-11-07 13:05:54.604167+00	2025-11-07 13:05:54.604167+00	\N
6754bb1e-d764-43c6-8e63-9347dcfbbfec	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	food	10000	monthly	\N	\N	f	80	t		2025-11-07 13:06:03.740615+00	2025-11-07 13:06:03.740615+00	\N
07df20f4-a407-4a7f-bada-918c5d04d3fb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	donation	5000	monthly	\N	\N	f	80	t		2025-11-07 13:06:13.974766+00	2025-11-07 13:06:13.974766+00	\N
\.


--
-- Data for Name: credit_card_payments; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.credit_card_payments (id, user_id, card_id, account_id, amount, payment_date, description, transaction_id, created_at, updated_at, deleted_at) FROM stdin;
f07125fd-9bd2-43e7-ba3a-d1d2af141cab	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	6079ac1c-7a40-45b0-913d-4e2dce11dab6	1614.81	2025-11-08 00:00:00+00	Credit card payment - Brac Bank (Credit Card)	ad1ac008-a63b-4ecc-999b-bbe1ad51e3c2	2025-11-08 18:04:44.458224+00	2025-11-08 18:04:44.458224+00	\N
d54a92d0-35c4-47fa-abb4-c3dc6de3a344	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	6079ac1c-7a40-45b0-913d-4e2dce11dab6	3898.23	2025-11-24 00:00:00+00	Credit card payment - Brac Bank (Credit Card)	a0177592-47be-49c2-80d5-962b6b95de69	2025-11-24 16:03:22.243487+00	2025-11-24 16:03:22.243487+00	\N
\.


--
-- Data for Name: credit_card_transactions; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.credit_card_transactions (id, user_id, card_id, transaction_id, category_id, amount, description, merchant, date, type, tags, attachments, created_at, updated_at, deleted_at) FROM stdin;
ebc99080-bd0d-4dd0-be3a-42a1534a4639	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	00000000-0000-0000-0000-000000000000		1614.81	Credit card payment - Brac Bank (Credit Card)		2025-11-08 00:00:00+00	payment	\N	\N	2025-11-08 18:04:44.461511+00	2025-11-08 18:04:44.461511+00	\N
e6b6c58b-b3c3-49ec-844b-7cba0e3c0b92	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	00000000-0000-0000-0000-000000000000		3898.23	Credit card payment - Brac Bank (Credit Card)		2025-11-24 00:00:00+00	payment	\N	\N	2025-11-24 16:03:22.250254+00	2025-11-24 16:03:22.250254+00	\N
\.


--
-- Data for Name: credit_cards; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.credit_cards (id, user_id, name, last_four_digits, card_network, credit_limit, current_balance, apr, due_date, statement_date, minimum_payment, last_payment_date, last_payment_amount, rewards_program, active, notes, created_at, updated_at, deleted_at) FROM stdin;
a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Brac Bank (Credit Card)	5025	Visa	200000	0	0	\N	\N	0	2025-11-24 00:00:00+00	3898.23		t		2025-11-01 17:00:54.255539+00	2025-11-24 16:03:22.243264+00	\N
\.


--
-- Data for Name: debt_payments; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.debt_payments (id, user_id, debt_id, account_id, amount, payment_date, description, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: debt_records; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.debt_records (id, user_id, creditor_name, original_amount, remaining_amount, account_id, status, borrowed_date, due_date, interest_rate, description, is_initial, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: goal_contributions; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.goal_contributions (id, user_id, goal_id, holding_id, type, amount, date, notes, transaction_id, created_at, updated_at, deleted_at) FROM stdin;
c333b630-815c-49a2-9e39-c15388f290ad	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	56ff16a0-de40-4241-8904-f6680afb2a61	57e0a724-872f-4f11-a6b5-72b946606ca6	contribution	13000	2024-09-12 10:37:00+00	External holding: DPS (Deposit Pension Scheme)	e62ea004-113d-4508-9ec2-21ee2c926bd5	2025-11-02 16:38:02.101663+00	2025-11-02 16:38:02.101663+00	\N
5cfc8a12-fc1b-48af-a920-e15f883e2b40	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	56ff16a0-de40-4241-8904-f6680afb2a61	33283ce5-6dff-4ede-99cd-37335e0870ff	contribution	1000	2025-11-13 08:43:00+00	Added DPS (Deposit Pension Scheme)	467b0679-ab62-4e84-a378-f7f2a9b0b54c	2025-11-13 14:44:10.4015+00	2025-11-13 14:44:10.4015+00	\N
8ed964d3-cf73-4e10-83f7-7ecad184dc16	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	33e31fd1-12ae-4df7-b31c-fda997d4e4ce	7ffa41df-6f62-4757-a60c-36ecd62066d8	contribution	55000	2024-05-18 10:43:00+00	External holding: DPS (Deposit Pension Scheme)	6b612542-991f-40c7-94d5-2c0ab3756034	2025-11-02 16:43:44.832174+00	2025-11-02 16:43:44.832174+00	\N
c65fd713-4164-4a74-bfa4-1770b016309b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	df2b69ff-147f-4676-af50-25db9ff1ee54	6fa890f0-f229-43f3-bc26-b0e9336c2ad5	contribution	4000	2025-07-24 10:06:00+00	External holding: DPS (Deposit Pension Scheme)	bba48dc4-2b64-4a1d-9bed-4e8b9ffeaaf3	2025-11-02 16:08:00.233478+00	2025-11-02 16:08:00.233478+00	\N
94822380-f44d-4f54-8c99-4276d26f64d1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	df2b69ff-147f-4676-af50-25db9ff1ee54	8c6b460b-97c0-4dc9-b19b-da2ae5c9e73b	contribution	1000	2025-11-01 10:46:00+00	Added DPS (Deposit Pension Scheme)	faee24b9-b925-48a0-87f4-268541d59b03	2025-11-02 16:47:01.999189+00	2025-11-02 16:47:01.999189+00	\N
15468af7-5b85-4a4b-8178-7c1ec162b40d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e0dc0db0-05c6-4ace-9f5d-6755859463db	72aea14c-ba23-45db-82a4-396dfc52b98b	contribution	50000	2026-06-09 10:20:00+00	External holding: DPS (Deposit Pension Scheme)	1f91f320-bac8-4739-897b-0dfeed161183	2025-11-02 16:21:46.399503+00	2025-11-02 16:21:46.399503+00	\N
9f7a5e27-1265-4a1b-8000-d2e3301788dc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e0dc0db0-05c6-4ace-9f5d-6755859463db	aa28dbe9-c6a3-413a-a780-c7fae3b4a675	contribution	10000	2025-11-01 10:47:00+00	Added DPS (Deposit Pension Scheme)	c6ee0331-3670-4842-a99b-1a5cbe41b92f	2025-11-02 16:47:33.518376+00	2025-11-02 16:47:33.518376+00	\N
824a0346-7434-451b-a6d7-e5de809572c8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e6b69c0c-d123-4f9b-a773-c818a5dfdaa9	1af1fc17-1b66-467a-8d74-dfae44329b6e	contribution	10000	2025-01-31 10:23:00+00	External holding: DPS (Deposit Pension Scheme)	fe8cf6ae-5ead-4ccf-a551-10289a99324b	2025-11-02 16:24:02.665788+00	2025-11-02 16:24:02.665788+00	\N
f1af8b5d-620d-4908-8ad6-fc69e2ef28cf	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e6b69c0c-d123-4f9b-a773-c818a5dfdaa9	b9db0210-9a3c-4331-bf3d-a3713d162d0e	contribution	10000	2025-11-01 10:47:00+00	Added DPS (Deposit Pension Scheme)	2e95fe7b-1114-4a03-bb4e-c3c2d526b741	2025-11-02 16:48:03.749461+00	2025-11-02 16:48:03.749461+00	\N
a3da3eab-b9da-4df2-beff-cc75034362b6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	84eb448d-76dd-46c6-a4b7-9c9179f6753e	888fa6cc-1842-4e34-91f5-95ef2c2bbf95	contribution	10000	2025-11-02 10:34:00+00	External holding: DPS (Deposit Pension Scheme)	1cbdee6a-62bd-48b2-bf4d-b99f2984c1f6	2025-11-02 16:35:03.543597+00	2025-11-02 16:35:03.543597+00	\N
3535e7bc-875c-4611-b720-92ff6e74efcb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	84eb448d-76dd-46c6-a4b7-9c9179f6753e	7a9c3ed7-8302-4b5d-8ac2-23e5c09f7734	contribution	1000	2025-11-01 10:50:00+00	Added DPS (Deposit Pension Scheme)	98ab9145-1d9b-4258-991c-4ec6c145ea68	2025-11-02 16:50:33.529803+00	2025-11-02 16:50:33.529803+00	\N
2f62e931-46c4-4f86-a215-f912fbadfd21	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	33e31fd1-12ae-4df7-b31c-fda997d4e4ce	28655439-8c56-4ed1-82b6-b0db265301e1	contribution	3000	2025-11-18 10:02:00+00	Added DPS (Deposit Pension Scheme)	eba09eb7-5ab3-48db-a5c7-6cadf68276ad	2025-11-18 16:03:01.055104+00	2025-11-18 16:03:01.055104+00	\N
3dd593b5-264f-41ee-9268-1e9e2e6b88fc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	2946b607-6a79-4617-8c2b-bbebd47657e5	c9d04081-a6ba-4b98-af44-d47ac05fcd2d	contribution	392000	2024-05-16 11:14:00+00	External holding: Fixed Deposit	24e09bac-0410-4256-b9ff-2136ed55feff	2025-11-02 17:16:55.900503+00	2025-11-02 17:16:55.900503+00	\N
ae927243-2c49-49d5-8f4b-3a2c70deff10	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	2946b607-6a79-4617-8c2b-bbebd47657e5	bf2b377f-442f-4d0b-93dc-6b0e61ecc0c4	contribution	12000	2025-11-01 11:17:00+00	Added Fixed Deposit	c4ca3200-5561-448f-ad23-56546fe3fc52	2025-11-02 17:17:16.924619+00	2025-11-02 17:17:16.924619+00	\N
19a87f3e-195f-4ad3-88d3-8609816aeee6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	dd8540c6-c548-4989-aa0f-41b1b379bafb	fe023801-4468-4e93-a81d-74f2236e4c4c	contribution	725000	2023-06-25 11:09:00+00	External holding: Fixed Deposit	5ba2b13e-0ec4-4caf-a755-8e716f9a3068	2025-11-02 17:10:37.265018+00	2025-11-02 17:10:37.265018+00	\N
ca31a851-d11f-4c69-a6ba-fb451e801a0f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	dd8540c6-c548-4989-aa0f-41b1b379bafb	fee052b0-82d8-4a34-ad3b-7fc265cd3a7d	contribution	25000	2025-11-01 11:10:00+00	Added Fixed Deposit	2d74ddb2-48a7-493a-bbf8-caf8d029f251	2025-11-02 17:11:10.00668+00	2025-11-02 17:11:10.00668+00	\N
6ea06fe7-0bf4-405a-8abc-4da5315a4ec6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6ab6926d-8ede-41c0-8cb0-a2f36ce6ea11	3f1231b8-4892-4e60-8996-90220b75a4f2	contribution	200000	2025-10-22 12:16:00+00	External holding: Fixed Deposit	05000557-7353-4e92-8c70-971e8912a77c	2025-11-01 18:17:53.866524+00	2025-11-01 18:17:53.866524+00	\N
\.


--
-- Data for Name: goal_holdings; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.goal_holdings (id, user_id, goal_id, name, type, status, purchase_date, amount, current_value, institution, account_number, interest_rate, maturity_date, maturity_amount, tenure_months, symbol, quantity, cost_basis, current_price, monthly_deposit, details, transaction_id, created_at, updated_at, deleted_at) FROM stdin;
1af1fc17-1b66-467a-8d74-dfae44329b6e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e6b69c0c-d123-4f9b-a773-c818a5dfdaa9	DPS (Deposit Pension Scheme)	dps	active	2025-01-31 10:23:00+00	10000	100000	bKash	\N	\N	2026-01-31 00:00:00+00	125593.47	\N		\N	\N	\N	10000	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:24:02.664724+00	2025-11-02 16:24:02.664724+00	\N
b9db0210-9a3c-4331-bf3d-a3713d162d0e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e6b69c0c-d123-4f9b-a773-c818a5dfdaa9	DPS (Deposit Pension Scheme)	savings	active	2025-11-01 10:47:00+00	10000	10000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:48:03.748123+00	2025-11-02 16:48:03.748123+00	\N
888fa6cc-1842-4e34-91f5-95ef2c2bbf95	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	84eb448d-76dd-46c6-a4b7-9c9179f6753e	DPS (Deposit Pension Scheme)	dps	active	2025-11-02 10:34:00+00	10000	10000	bKash	\N	\N	2026-01-09 00:00:00+00	12559.33	\N		\N	\N	\N	999.97	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:35:03.542715+00	2025-11-02 16:35:03.542715+00	\N
7a9c3ed7-8302-4b5d-8ac2-23e5c09f7734	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	84eb448d-76dd-46c6-a4b7-9c9179f6753e	DPS (Deposit Pension Scheme)	savings	active	2025-11-01 10:50:00+00	1000	1000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:50:33.528709+00	2025-11-02 16:50:33.528709+00	\N
57e0a724-872f-4f11-a6b5-72b946606ca6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	56ff16a0-de40-4241-8904-f6680afb2a61	DPS (Deposit Pension Scheme)	dps	active	2024-09-12 10:37:00+00	13000	13000	bKash	\N	\N	2026-09-12 00:00:00+00	26340	\N		\N	\N	\N	1000	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:38:02.100837+00	2025-11-02 16:38:02.100837+00	\N
33283ce5-6dff-4ede-99cd-37335e0870ff	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	56ff16a0-de40-4241-8904-f6680afb2a61	DPS (Deposit Pension Scheme)	savings	active	2025-11-13 08:43:00+00	1000	1000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-13 14:44:10.398544+00	2025-11-13 14:44:10.398544+00	\N
7ffa41df-6f62-4757-a60c-36ecd62066d8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	33e31fd1-12ae-4df7-b31c-fda997d4e4ce	DPS (Deposit Pension Scheme)	dps	active	2024-05-18 10:43:00+00	55000	55000	bKash	\N	\N	2026-05-18 00:00:00+00	77775	\N		\N	\N	\N	3000	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:43:44.83132+00	2025-11-02 16:43:44.83132+00	\N
28655439-8c56-4ed1-82b6-b0db265301e1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	33e31fd1-12ae-4df7-b31c-fda997d4e4ce	DPS (Deposit Pension Scheme)	savings	active	2025-11-18 10:02:00+00	3000	3000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-18 16:03:01.053203+00	2025-11-18 16:03:01.053203+00	\N
c9d04081-a6ba-4b98-af44-d47ac05fcd2d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	2946b607-6a79-4617-8c2b-bbebd47657e5	Fixed Deposit	fixed_deposit	active	2024-05-16 11:14:00+00	392000	392000	IDLC	\N	\N	2029-06-01 00:00:00+00	\N	60		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 17:16:55.898815+00	2025-11-02 17:16:55.898815+00	\N
bf2b377f-442f-4d0b-93dc-6b0e61ecc0c4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	2946b607-6a79-4617-8c2b-bbebd47657e5	Fixed Deposit	savings	active	2025-11-01 11:17:00+00	12000	12000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 17:17:16.923619+00	2025-11-02 17:17:16.923619+00	\N
fee052b0-82d8-4a34-ad3b-7fc265cd3a7d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	dd8540c6-c548-4989-aa0f-41b1b379bafb	Fixed Deposit	savings	active	2025-11-01 11:10:00+00	25000	25000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 17:11:10.00478+00	2025-11-02 17:11:10.00478+00	\N
fe023801-4468-4e93-a81d-74f2236e4c4c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	dd8540c6-c548-4989-aa0f-41b1b379bafb	Fixed Deposit	fixed_deposit	active	2023-06-25 05:09:00+00	725000	1025000	IDLC	\N	\N	2030-11-02 00:00:00+00	\N	60		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 17:10:37.263633+00	2025-11-02 17:12:38.938691+00	\N
3f1231b8-4892-4e60-8996-90220b75a4f2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6ab6926d-8ede-41c0-8cb0-a2f36ce6ea11	Fixed Deposit	fixed_deposit	active	2025-10-22 12:16:00+00	200000	200000	Brac Bank	\N	9.5	2026-11-22 00:00:00+00	220900	13		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-01 18:17:53.858707+00	2025-11-01 18:17:53.858707+00	\N
6fa890f0-f229-43f3-bc26-b0e9336c2ad5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	df2b69ff-147f-4676-af50-25db9ff1ee54	DPS (Deposit Pension Scheme)	dps	active	2025-07-24 10:06:00+00	4000	4000	bKash	\N	\N	2026-07-24 00:00:00+00	12597.25	\N		\N	\N	\N	1000	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:08:00.232469+00	2025-11-02 16:08:00.232469+00	\N
8c6b460b-97c0-4dc9-b19b-da2ae5c9e73b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	df2b69ff-147f-4676-af50-25db9ff1ee54	DPS (Deposit Pension Scheme)	savings	active	2025-11-01 10:46:00+00	1000	1000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:47:01.997953+00	2025-11-02 16:47:01.997953+00	\N
72aea14c-ba23-45db-82a4-396dfc52b98b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e0dc0db0-05c6-4ace-9f5d-6755859463db	DPS (Deposit Pension Scheme)	dps	active	2026-06-09 10:20:00+00	50000	50000	bKash	\N	\N	2026-06-09 00:00:00+00	125972.5	\N		\N	\N	\N	10000	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:21:46.398318+00	2025-11-02 16:21:46.398318+00	\N
aa28dbe9-c6a3-413a-a780-c7fae3b4a675	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e0dc0db0-05c6-4ace-9f5d-6755859463db	DPS (Deposit Pension Scheme)	savings	active	2025-11-01 10:47:00+00	10000	10000		\N	\N	\N	\N	\N		\N	\N	\N	\N	\N	00000000-0000-0000-0000-000000000000	2025-11-02 16:47:33.517223+00	2025-11-02 16:47:33.517223+00	\N
\.


--
-- Data for Name: goals; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.goals (id, user_id, name, description, icon, color, category, priority, target_amount, current_amount, target_date, monthly_contribution, status, achieved, achieved_date, last_contribution, last_contribution_date, created_at, updated_at, deleted_at) FROM stdin;
6c16099e-7576-4170-9488-c2a068dc7c12	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash DPS - 2145230221980		🎯	#3b82f6	other	high	120000	0	2026-01-09 00:00:00+00	10000	active	f	\N	0	\N	2025-11-02 16:32:22.472492+00	2025-11-02 16:32:58.31222+00	2025-11-02 16:33:00.685571+00
d91cf9e1-2d1c-453f-9e4c-19aac38d8679	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Brac Bank Fixed Deposit	Opening on 22 Nov 2025 at a 9.50% interest rate	🎯	#3b82f6	other	high	200000	0	2026-11-22 00:00:00+00	0	active	f	\N	0	\N	2025-11-01 17:03:38.944128+00	2025-11-01 17:05:34.497706+00	2025-11-01 17:05:37.149507+00
df2b69ff-147f-4676-af50-25db9ff1ee54	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash DPS - 2145230866410		🎯	#3b82f6	other	high	12000	5000	2026-07-24 00:00:00+00	1000	active	f	\N	1000	2025-11-01 10:46:00+00	2025-11-02 16:06:52.837283+00	2025-11-29 07:43:19.794406+00	\N
e0dc0db0-05c6-4ace-9f5d-6755859463db	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash DPS - 2145230701614		🎯	#3b82f6	other	high	120000	60000	2026-06-09 00:00:00+00	10000	active	f	\N	10000	2025-11-01 10:47:00+00	2025-11-02 16:20:40.713877+00	2025-11-29 07:43:19.798128+00	\N
e6b69c0c-d123-4f9b-a773-c818a5dfdaa9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash DPS - 2145230291760		🎯	#3b82f6	other	high	120000	110000	2026-01-31 00:00:00+00	10000	active	f	\N	10000	2025-11-01 10:47:00+00	2025-11-02 16:23:03.244867+00	2025-11-29 07:43:19.800696+00	\N
84eb448d-76dd-46c6-a4b7-9c9179f6753e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash DPS - 2145230221980		🎯	#3b82f6	other	high	12000	11000	2026-01-09 00:00:00+00	1000	active	f	\N	1000	2025-11-01 10:50:00+00	2025-11-02 16:34:21.911244+00	2025-11-29 07:43:19.803719+00	\N
56ff16a0-de40-4241-8904-f6680afb2a61	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash DPS - 1783060326999		🎯	#3b82f6	other	high	24000	14000	2026-09-12 00:00:00+00	1000	active	f	\N	1000	2025-11-13 08:43:00+00	2025-11-02 16:37:17.7819+00	2025-11-29 07:43:19.806174+00	\N
33e31fd1-12ae-4df7-b31c-fda997d4e4ce	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	bKash DPS - 1783060302607		🎯	#3b82f6	other	high	72000	58000	2026-05-18 00:00:00+00	3000	active	f	\N	3000	2025-11-18 10:02:00+00	2025-11-02 16:43:04.930308+00	2025-11-29 07:43:19.809102+00	\N
2946b607-6a79-4617-8c2b-bbebd47657e5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	IDLC ISF-SIP-001807		🎯	#3b82f6	other	high	720000	404000	2029-05-01 00:00:00+00	12000	active	f	\N	12000	2025-11-01 11:17:00+00	2025-11-02 17:14:11.019772+00	2025-11-29 07:43:19.815155+00	\N
dd8540c6-c548-4989-aa0f-41b1b379bafb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	IDLC ISF-SIP-001449		🎯	#3b82f6	other	high	1500000	1050000	2028-06-01 00:00:00+00	25000	active	f	\N	25000	2025-11-01 11:10:00+00	2025-11-02 17:09:18.477154+00	2025-11-29 07:43:19.81783+00	\N
6ab6926d-8ede-41c0-8cb0-a2f36ce6ea11	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Brac Bank Fixed Deposit	Invested 200000 on 9.50% interest rate on 22 October 2025	🎯	#3b82f6	other	high	200000	200000	2026-11-22 00:00:00+00	0	active	f	\N	200000	2025-10-22 12:16:00+00	2025-11-01 18:16:49.105765+00	2025-11-29 07:43:19.820568+00	\N
\.


--
-- Data for Name: lend_payments; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.lend_payments (id, user_id, lend_id, account_id, amount, payment_date, description, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: lend_records; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.lend_records (id, user_id, debtor_name, original_amount, remaining_amount, account_id, status, lent_date, due_date, interest_rate, description, is_initial, created_at, updated_at, deleted_at) FROM stdin;
0f34ad12-adea-40f5-ac80-0cde89b0c630	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Sohug Mama	140000	140000	\N	active	2024-01-01 00:00:00	\N	\N		f	2025-11-08 08:25:09.904137+00	2025-11-08 08:25:09.904137+00	\N
55821a5e-b9f0-4cb4-a60f-084f5bd76f44	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Ripon Bhai	277500	277500	\N	active	2023-01-01 00:00:00	\N	\N		f	2025-11-08 18:07:36.39057+00	2025-11-08 18:07:36.39057+00	\N
9bf5b5ac-6514-41e0-bdb3-5328c6970dd1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Nurul Huda	2500	2500	c3288300-472e-413e-9588-48542b2f66d6	active	2025-11-21 00:00:00	\N	\N		f	2025-11-22 16:51:06.153112+00	2025-11-22 16:51:06.153112+00	\N
a19e8944-a810-4272-a46e-2fe47fdcb35f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	Jamal bhai	2037	2037	c3288300-472e-413e-9588-48542b2f66d6	active	2025-11-28 00:00:00	\N	\N		f	2025-11-29 03:25:54.112293+00	2025-11-29 03:25:54.112293+00	\N
\.


--
-- Data for Name: reconciliation_transactions; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.reconciliation_transactions (id, reconciliation_id, transaction_id, created_at) FROM stdin;
\.


--
-- Data for Name: reconciliations; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.reconciliations (id, user_id, account_id, reconciliation_date, statement_balance, book_balance, difference, notes, status, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: recurring_transactions; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.recurring_transactions (id, user_id, template_id, template_user_id, template_account_id, template_to_account_id, template_type, template_amount, template_category_id, template_date, template_description, template_tags, template_savings_goal_id, template_fixed_deposit_id, template_investment_id, template_recurring_id, template_credit_card_id, template_attachments, template_reconciled, template_reconciliation_id, template_created_at, template_updated_at, template_deleted_at, frequency, start_date, end_date, last_processed, enabled, created_at, updated_at, deleted_at) FROM stdin;
8845435a-f95f-448b-850a-dfb53371cf08	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	834d03e7-7c59-4b7a-98ce-7c35a392ba47	00000000-0000-0000-0000-000000000000	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	0001-01-01 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	\N	\N	\N	f	\N	2025-11-07 13:00:54.222883+00	2025-11-15 08:02:35.120555+00	2025-11-15 08:03:08.149458+00	daily	2025-11-08 06:00:00	2027-01-01 06:00:00	2025-11-15 08:02:35.106628+00	t	2025-11-07 13:00:54.222883+00	2025-11-15 08:02:35.120572+00	\N
8845435a-f95f-448b-850a-dfb53371cf08	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	e400ab80-bdbc-44b2-9bf9-68ce033ff248	00000000-0000-0000-0000-000000000000	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	0001-01-01 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	\N	\N	\N	f	\N	2025-11-15 08:02:33.15156+00	2025-11-15 08:02:35.121245+00	2025-11-15 08:03:14.665661+00	daily	2025-11-15 00:00:00	2027-01-01 00:00:00	2025-11-15 08:02:35.106628+00	t	2025-11-07 13:00:54.222883+00	2025-11-15 08:02:35.121248+00	\N
\.


--
-- Data for Name: rewards; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.rewards (id, user_id, card_id, type, amount, description, earned_date, redeemed, redeemed_at, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: service_records; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.service_records (id, user_id, asset_id, service_date, service_type, service_provider, cost, description, notes, warranty_covered, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: settings; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.settings (id, user_id, currency, dark_mode, date_format, first_day_of_week, language, notif_push, notif_email, notif_budget_alerts, notif_bill_reminders, created_at, updated_at, deleted_at) FROM stdin;
50301d40-6153-414c-be16-79f58947269d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	BDT	f	MM/DD/YYYY	0	en	t	t	t	t	2025-11-01 11:57:04.987412+00	2025-11-02 17:38:36.646262+00	\N
\.


--
-- Data for Name: statements; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.statements (id, user_id, card_id, statement_date, due_date, opening_balance, closing_balance, minimum_payment, total_charges, total_payments, interest_charged, paid, paid_date, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.tags (id, user_id, name, color, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.transactions (id, user_id, account_id, to_account_id, type, amount, category_id, date, description, tags, savings_goal_id, fixed_deposit_id, investment_id, recurring_id, credit_card_id, attachments, reconciled, reconciliation_id, created_at, updated_at, deleted_at) FROM stdin;
20927be9-b4f8-4c1a-bd25-4a61bc1fbb8a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	income	24867.13	opening_balance	2025-11-01 12:41:34.973+00	Opening balance for bKash	\N	\N	\N	\N	\N	\N	\N	f	\N	2025-11-01 12:41:34.977901+00	2025-11-01 12:41:34.977901+00	\N
2358ced5-3577-4d62-9fd6-015de3d01686	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	income	151650.92	opening_balance	2025-11-01 12:42:03.199886+00	Opening balance for SCB Debit	\N	\N	\N	\N	\N	\N	\N	f	\N	2025-11-01 12:42:03.200322+00	2025-11-01 12:42:03.200322+00	\N
68a47ac4-8087-47b1-ae4c-410628e07c2f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	\N	income	355727.12	opening_balance	2025-11-01 16:59:03.851641+00	Opening balance for Brac Bank Debit	\N	\N	\N	\N	\N	\N	\N	f	\N	2025-11-01 16:59:03.852271+00	2025-11-01 16:59:03.852271+00	\N
05000557-7353-4e92-8c70-971e8912a77c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	200000	goal_external_holding	2025-10-22 12:16:00+00	External Fixed Deposit tracked for Brac Bank Fixed Deposit	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-01 18:17:53.864547+00	2025-11-01 18:17:53.864547+00	\N
7c505ec8-bdb2-4898-a6fe-7d1d8f30b3ff	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-24 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-24 03:41:55.123803+00	2025-11-24 03:41:55.123803+00	\N
bba48dc4-2b64-4a1d-9bed-4e8b9ffeaaf3	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	4000	goal_external_holding	2025-07-24 10:06:00+00	External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230866410	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:08:00.233165+00	2025-11-02 16:08:00.233165+00	\N
1f91f320-bac8-4739-897b-0dfeed161183	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	50000	goal_external_holding	2026-06-09 10:20:00+00	External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230701614	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:21:46.399021+00	2025-11-02 16:21:46.399021+00	\N
fe8cf6ae-5ead-4ccf-a551-10289a99324b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	10000	goal_external_holding	2025-01-31 10:23:00+00	External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230291760	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:24:02.665361+00	2025-11-02 16:24:02.665361+00	\N
1cbdee6a-62bd-48b2-bf4d-b99f2984c1f6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	10000	goal_external_holding	2025-11-02 10:34:00+00	External DPS (Deposit Pension Scheme) tracked for bKash DPS - 2145230221980	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:35:03.543314+00	2025-11-02 16:35:03.543314+00	\N
e62ea004-113d-4508-9ec2-21ee2c926bd5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	13000	goal_external_holding	2024-09-12 10:37:00+00	External DPS (Deposit Pension Scheme) tracked for bKash DPS - 1783060326999	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:38:02.101328+00	2025-11-02 16:38:02.101328+00	\N
6b612542-991f-40c7-94d5-2c0ab3756034	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	55000	goal_external_holding	2024-05-18 10:43:00+00	External DPS (Deposit Pension Scheme) tracked for bKash DPS - 1783060302607	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:43:44.831818+00	2025-11-02 16:43:44.831818+00	\N
faee24b9-b925-48a0-87f4-268541d59b03	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1000	goal_holding_added	2025-11-01 10:46:00+00	Added to bKash DPS - 2145230866410: DPS (Deposit Pension Scheme)	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:47:01.998419+00	2025-11-02 16:47:01.998419+00	\N
c6ee0331-3670-4842-a99b-1a5cbe41b92f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	10000	goal_holding_added	2025-11-01 10:47:00+00	Added to bKash DPS - 2145230701614: DPS (Deposit Pension Scheme)	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:47:33.517648+00	2025-11-02 16:47:33.517648+00	\N
2e95fe7b-1114-4a03-bb4e-c3c2d526b741	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	10000	goal_holding_added	2025-11-01 10:47:00+00	Added to bKash DPS - 2145230291760: DPS (Deposit Pension Scheme)	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:48:03.748527+00	2025-11-02 16:48:03.748527+00	\N
98ab9145-1d9b-4258-991c-4ec6c145ea68	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1000	goal_holding_added	2025-11-01 10:50:00+00	Added to bKash DPS - 2145230221980: DPS (Deposit Pension Scheme)	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 16:50:33.529188+00	2025-11-02 16:50:33.529188+00	\N
5ba2b13e-0ec4-4caf-a755-8e716f9a3068	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	725000	goal_external_holding	2023-06-25 11:09:00+00	External Fixed Deposit tracked for IDLC ISF-SIP-001449	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 17:10:37.264537+00	2025-11-02 17:10:37.264537+00	\N
2d74ddb2-48a7-493a-bbf8-caf8d029f251	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	expense	25000	goal_holding_added	2025-11-01 11:10:00+00	Added to IDLC ISF-SIP-001449: Fixed Deposit	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 17:11:10.005217+00	2025-11-02 17:11:10.005217+00	\N
24e09bac-0410-4256-b9ff-2136ed55feff	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	tracking	392000	goal_external_holding	2024-05-16 11:14:00+00	External Fixed Deposit tracked for IDLC ISF-SIP-001807	["goal", "holding", "external", "tracking", "hidden"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 17:16:55.899574+00	2025-11-02 17:16:55.899574+00	\N
c4ca3200-5561-448f-ad23-56546fe3fc52	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	expense	12000	goal_holding_added	2025-11-01 11:17:00+00	Added to IDLC ISF-SIP-001807: Fixed Deposit	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 17:17:16.924063+00	2025-11-02 17:17:16.924063+00	\N
ef95c42e-a4bf-4330-9fc9-431bd67f94e8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	income	50440	opening_balance	2025-11-02 17:37:12.975949+00	Opening balance for Money Bag	\N	\N	\N	\N	\N	\N	\N	f	\N	2025-11-02 17:37:12.978036+00	2025-11-02 17:37:12.978036+00	\N
eaaf61a5-c5fe-47f8-a2f3-260d2b43a151	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-24 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-24 15:00:53.047124+00	2025-11-24 15:00:53.047124+00	\N
4e618279-5e38-47b3-91a2-e2f111270099	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-24 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-24 15:01:12.495391+00	2025-11-24 15:01:12.495391+00	\N
02a535dd-be27-41dd-8bc3-4991f9886ac1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	500	donation	2025-11-24 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-24 15:01:24.985811+00	2025-11-24 15:01:24.985811+00	\N
dfb32f1c-9d8c-4376-9f54-27362e786b9a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	20	donation	2025-11-23 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-24 15:16:03.073299+00	2025-11-24 15:16:03.073299+00	\N
12adafa0-f4b9-4ad3-85f0-83daa03477dc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	expense	1745.99	other_expense	2025-11-22 00:00:00+00	Mess bajar	[]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	[]	f	\N	2025-11-24 15:57:15.018772+00	2025-11-24 15:57:15.018772+00	\N
a56ce51d-07a7-4f8b-b109-94f995380011	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	expense	609.3	other_expense	2025-11-23 00:00:00+00	Mess Bazer	[]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	[]	f	\N	2025-11-24 15:57:48.250503+00	2025-11-24 15:57:48.250503+00	\N
81bb1c6f-4850-45ba-a793-0f9a0f8859bc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	expense	466	food	2025-11-23 00:00:00+00	Fruits	[]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	[]	f	\N	2025-11-24 15:58:07.802648+00	2025-11-24 15:58:07.802648+00	\N
692eee26-d26d-4dd9-9207-70982b6ac8ea	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	expense	796.94	other_expense	2025-11-20 00:00:00+00	Mess Bazar	[]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	[]	f	\N	2025-11-24 16:01:51.475336+00	2025-11-24 16:01:51.475336+00	\N
36372654-9013-40dd-960c-ee9eb5c84f58	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	expense	280	food	2025-11-20 00:00:00+00	Fruits	[]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	[]	f	\N	2025-11-24 16:02:10.4532+00	2025-11-24 16:02:10.4532+00	\N
a0177592-47be-49c2-80d5-962b6b95de69	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	\N	expense	3898.23	credit_card_payment	2025-11-24 00:00:00+00	Credit card payment - Brac Bank (Credit Card)	["credit_card_payment"]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	\N	f	\N	2025-11-24 16:03:22.242569+00	2025-11-24 16:03:22.242569+00	\N
285d505a-fb8c-42da-ac0a-4cf79b06af3c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	110	transport	2025-11-25 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-25 15:52:34.356119+00	2025-11-25 15:52:34.356119+00	\N
c9f7943a-69ee-4461-86e5-83fe8f5b18c5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-25 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-25 15:52:47.19733+00	2025-11-25 15:52:47.19733+00	\N
42024371-e47e-4464-ae20-4b759566c5c0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-25 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-25 15:55:35.694034+00	2025-11-25 15:55:35.694034+00	\N
7874feb8-44eb-42e5-a3f4-1ee9f99630ce	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	income	137283	salary	2025-11-25 00:00:00+00	Nov 2025/Salary	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-25 15:56:13.764187+00	2025-11-25 15:56:13.764187+00	\N
7d695631-2ebb-4aec-985d-77078cb350cf	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	income	14500	salary	2025-11-25 00:00:00+00	Nov 2025/Salary	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-25 15:56:30.593236+00	2025-11-25 15:56:30.593236+00	\N
3a60ec7f-30ff-452a-aa16-5c3b96c0ebf2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	income	100	other_income	2025-11-26 00:00:00+00	From Lukman bhai for Youtube	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:30:36.314003+00	2025-11-29 03:30:36.314003+00	\N
72e1f60f-fb50-4ff8-8f3d-b01b6e058bd6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	income	100	other_income	2025-11-29 00:00:00+00	Cashback	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:37:05.407735+00	2025-11-29 03:37:05.407735+00	\N
467b0679-ab62-4e84-a378-f7f2a9b0b54c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1000	goal_holding_added	2025-11-13 08:43:00+00	Added to bKash DPS - 1783060326999: DPS (Deposit Pension Scheme)	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-13 14:44:10.400913+00	2025-11-13 14:44:10.400913+00	\N
31a15c49-17e9-4e17-8036-77f7e6ef27a2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-26 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:29:54.971686+00	2025-11-27 16:29:54.971686+00	\N
4ca91ce0-6e8c-43de-b611-d3caaf07fe84	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	110	transport	2025-11-26 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:30:08.76301+00	2025-11-27 16:30:08.76301+00	\N
33cd2521-0b2b-4556-ad23-a7701ff7d747	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	110	transport	2025-11-27 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:30:18.658325+00	2025-11-27 16:30:18.658325+00	\N
f5948d5c-b2c7-46b2-a12e-78b18f635988	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	200	transport	2025-11-27 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:30:31.384112+00	2025-11-27 16:30:31.384112+00	\N
57ad2986-8eb1-4ced-a93d-2d38b2b07706	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	240	transport	2025-11-27 00:00:00+00	CNG fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:30:50.146562+00	2025-11-27 16:30:50.146562+00	\N
4964d3cf-b729-493e-b039-ede372efd9e1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	income	6274	other_income	2025-11-27 00:00:00+00	Insurance	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:31:47.311571+00	2025-11-27 16:31:47.311571+00	\N
f5088e9d-ce7e-496c-af7a-e9245b31caba	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1120	donation	2025-11-26 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:32:11.536755+00	2025-11-27 16:32:11.536755+00	\N
5456aa2f-1a77-4fd6-a289-7b0f7591a822	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-27 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:32:21.906759+00	2025-11-27 16:32:21.906759+00	\N
c45fd1ab-8363-463a-ae36-bd41d8af0406	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	500	donation	2025-11-27 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:32:31.731731+00	2025-11-27 16:32:31.731731+00	\N
3a43a967-62f2-4797-8b09-f71579a96537	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	380	other_expense	2025-11-27 00:00:00+00	Mess Bazar	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-27 16:33:23.079245+00	2025-11-27 16:33:23.079245+00	\N
b8d126d1-b9cc-4ef5-97fa-e8fdb9005476	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	140	transport	2025-11-28 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:10:01.428247+00	2025-11-29 03:10:01.428247+00	\N
52b975d5-1a20-4ef4-9c22-7b1350214e6a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	220	transport	2025-11-28 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:10:26.342301+00	2025-11-29 03:10:26.342301+00	\N
e45eb6d3-e49d-4ee4-8b07-fd8f49fc4d63	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	220	transport	2025-11-28 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:10:53.975306+00	2025-11-29 03:10:53.975306+00	\N
8f57d6b9-5eec-404e-a4a9-e8737d82363e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-28 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:16:31.370846+00	2025-11-29 03:16:31.370846+00	\N
1be42af4-7828-4338-aaf3-3292b87f1d9e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	60	donation	2025-11-28 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:17:02.339733+00	2025-11-29 03:17:02.339733+00	\N
83464a1c-f014-4bfc-99b2-1ae5a7c823c4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	500	donation	2025-11-28 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:18:19.062742+00	2025-11-29 03:18:19.062742+00	\N
a22ea27c-373f-4a19-b32d-483d75005f82	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	100	utilities	2025-11-28 00:00:00+00	Mobile recharge	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:18:50.047508+00	2025-11-29 03:18:50.047508+00	\N
8740f4f9-b89d-40dd-9ae0-144ac0a6a4e0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	765	healthcare	2025-11-28 00:00:00+00	Medicine	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:19:26.978268+00	2025-11-29 03:19:26.978268+00	\N
56b94022-7c6f-40bc-87e2-d88b54e1ded0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	9166.5	donation	2025-11-28 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:20:03.479202+00	2025-11-29 03:20:03.479202+00	\N
38b3e6d9-fc0e-4026-9b99-f55d3f33a163	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	30	food	2025-11-28 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:20:36.497916+00	2025-11-29 03:20:36.497916+00	\N
479aafc0-7179-4066-9273-ca129f956e26	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	505	food	2025-11-28 00:00:00+00	Raqi	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:21:40.247536+00	2025-11-29 03:21:40.247536+00	\N
0fbf67b0-0c14-48a4-adfe-0576c0193167	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	466.3	food	2025-11-28 00:00:00+00	From Foodpanda	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:22:48.032113+00	2025-11-29 03:22:48.032113+00	\N
3c763210-e1f2-4ee6-853a-4bb93f0ae105	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	649.99	other_expense	2025-11-28 00:00:00+00	For mess	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:23:30.285636+00	2025-11-29 03:23:30.285636+00	\N
33c98a2c-b880-4b85-beb4-e6d453283839	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	270	other_expense	2025-11-28 00:00:00+00	For mess	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:24:56.761515+00	2025-11-29 03:24:56.761515+00	\N
eba09eb7-5ab3-48db-a5c7-6cadf68276ad	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	3000	goal_holding_added	2025-11-18 10:02:00+00	Added to bKash DPS - 1783060302607: DPS (Deposit Pension Scheme)	["goal", "holding"]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-18 16:03:01.053897+00	2025-11-18 16:03:01.053897+00	\N
b28574d2-c2ab-4cae-9b04-cd9007826c4b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	2037	lend	2025-11-29 00:00:00+00	Lent to Jamal bhai	\N	\N	\N	\N	\N	\N	[]	f	\N	2025-11-29 03:25:54.116875+00	2025-11-29 03:30:50.468557+00	\N
1b9f8926-a5ba-4308-aa5c-45917b2149f7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-02 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-02 15:48:34.224044+00	2025-11-02 15:48:34.224044+00	\N
c3802ca2-dbc0-4652-8065-7782a3b2b1e9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	499	donation	2025-11-02 00:00:00+00	Donation for Nuts	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-02 15:49:50.059598+00	2025-11-02 15:49:50.059598+00	\N
a71302a1-0030-4284-b654-278088dd54bf	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1500	utilities	2025-11-01 00:00:00+00	Internet bill	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-02 16:51:05.285357+00	2025-11-02 16:51:05.285357+00	\N
fad607ef-b4f7-415b-804f-d1655b0fa695	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-01 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-02 16:51:44.236531+00	2025-11-02 16:51:44.236531+00	\N
96f5ce51-59e5-4257-af1d-40f45559b1e7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	500	utilities	2025-11-01 00:00:00+00	DESCO prepaid bill	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-02 16:52:39.859276+00	2025-11-02 16:52:39.859276+00	\N
c099b99b-3e6b-480c-8d5e-da9eb84a77b6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	140	transport	2025-11-02 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-02 17:37:31.666623+00	2025-11-02 17:37:31.666623+00	\N
a646668c-1f77-408e-bff8-165e480a4d89	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-02 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-02 17:37:52.604777+00	2025-11-02 17:37:52.604777+00	\N
c702dfba-6359-4b4f-9777-8e67e53c7829	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-03 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-03 05:20:14.088709+00	2025-11-03 05:20:14.088709+00	\N
abbafb8a-4521-4155-815e-82af3798f4e9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	130	housing	2025-11-03 00:00:00+00	Garbage Bag for Mess	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-03 15:53:42.270723+00	2025-11-03 15:53:42.270723+00	\N
eb273fef-13dd-419d-92f7-c425b6d0d721	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-03 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-03 15:54:01.348217+00	2025-11-03 15:54:01.348217+00	\N
acb6339d-1d2c-4257-9860-2c38264b7a4f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-03 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-03 15:54:40.473078+00	2025-11-03 15:54:40.473078+00	\N
3c02d554-581e-4c36-8b3a-f1bb37709218	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	c3288300-472e-413e-9588-48542b2f66d6	transfer	5000	transfer	2025-11-04 00:00:00+00	Transfer between accounts	[]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-04 13:31:16.59+00	2025-11-04 13:31:16.59+00	\N
065c91fe-30c3-4e06-abdc-3c0a76733a52	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1505	food	2025-11-04 00:00:00+00	Secret recipe	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-05 04:09:14.787251+00	2025-11-05 04:09:14.787251+00	\N
97baf158-9df0-4df0-8936-4e237c2a3250	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-05 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-05 04:09:35.761072+00	2025-11-05 04:09:35.761072+00	\N
6a66ed19-e95a-4e2e-b349-6be2125f5a5e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	10	donation	2025-11-05 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-05 04:09:58.578124+00	2025-11-05 04:09:58.578124+00	\N
eb8be643-5247-4672-a930-b6dd813cddf6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	200	donation	2025-11-04 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-05 04:10:46.330631+00	2025-11-05 04:10:46.330631+00	\N
5cad292a-0597-4ea2-976e-200ec3af2df7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-05 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-05 15:52:58.350899+00	2025-11-05 15:52:58.350899+00	\N
3affaafc-ec1d-4487-b823-d9315edc8007	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-05 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-05 15:54:50.720092+00	2025-11-05 15:54:50.720092+00	\N
0df44019-416c-4654-bdb8-ad5368f2cdc7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-06 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-07 08:42:47.814545+00	2025-11-07 08:42:47.814545+00	\N
e4976e8d-add7-4c47-9cfe-14d6a5c988e1	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-07 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-07 08:43:01.602677+00	2025-11-07 08:43:01.602677+00	\N
c5f5e973-f8d7-40c6-a473-304ba50e8e36	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	90	transport	2025-11-07 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-07 08:43:16.878293+00	2025-11-07 08:43:16.878293+00	\N
5791f38a-7e9a-4f11-a4d9-052b64b59c54	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	400	food	2025-11-06 00:00:00+00	Sarah resort pepsi	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-07 08:43:44.962968+00	2025-11-07 08:43:44.962968+00	\N
7f51ae39-3f79-4d8e-9ba9-7e1c1436e0ab	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	70	transport	2025-11-06 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-07 13:08:27.5759+00	2025-11-07 13:08:27.5759+00	\N
f04b6cf5-49cf-4c3b-a20d-9a2fa8a94dad	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-08 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-08 06:46:44.986574+00	2025-11-08 06:46:44.986574+00	\N
3d1b594e-f394-4600-9406-253e43b14621	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	598	utilities	2025-11-08 00:00:00+00	Mobile recharge for Maa	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-08 15:06:28.161034+00	2025-11-08 15:06:28.161034+00	\N
54f40500-5289-4cc6-b4fe-1068debf484a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1000	utilities	2025-11-13 00:00:00+00	DESCO Bill Pay	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-13 17:15:49.624402+00	2025-11-13 17:15:49.624402+00	\N
1ee7cc20-fa4d-4693-90c0-8ccd1cff5ce2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	expense	1014	food	2025-11-08 00:00:00+00	Food for Mess	[]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	[]	f	\N	2025-11-08 18:02:33.184802+00	2025-11-08 18:03:29.004866+00	\N
a959db34-3330-432d-8a8c-14e70b15fbcb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-14 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-14 11:28:31.654918+00	2025-11-14 11:28:31.654918+00	\N
d62f9d57-7874-4292-8497-d137e3a8a609	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	00000000-0000-0000-0000-000000000000	\N	expense	600.81	food	2025-11-08 00:00:00+00	Guava and Rice bran oil	[]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	[]	f	\N	2025-11-08 18:02:07.645428+00	2025-11-08 18:04:33.566765+00	\N
ad1ac008-a63b-4ecc-999b-bbe1ad51e3c2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	\N	expense	1614.81	credit_card_payment	2025-11-08 00:00:00+00	Credit card payment - Brac Bank (Credit Card)	["credit_card_payment"]	\N	\N	\N	\N	a2dc8b35-ccf6-4ff1-8be8-343d10d548ec	\N	f	\N	2025-11-08 18:04:44.455194+00	2025-11-08 18:04:44.455194+00	\N
c272af47-ed23-473a-b866-1ef1269d2aaa	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	114	donation	2025-11-09 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-09 16:58:14.781816+00	2025-11-09 16:58:14.781816+00	\N
9e4fad59-4d4b-4bf6-97d5-dfc1200ec81f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	200	transport	2025-11-09 00:00:00+00	CNG fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-09 16:58:35.862246+00	2025-11-09 16:58:35.862246+00	\N
87b3d61e-80ca-4e14-99fc-1d0763d98efd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-09 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-09 16:58:53.658473+00	2025-11-09 16:58:53.658473+00	\N
1394b9c5-a5bb-4d2d-b265-610a264a4c20	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	\N	expense	38802	other_expense	2025-11-09 00:00:00+00	TAX pay	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-09 16:59:45.642014+00	2025-11-09 16:59:45.642014+00	\N
e75f993c-d932-47f1-b47a-61d116364867	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	150	transport	2025-11-09 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-09 17:01:16.602348+00	2025-11-09 17:01:16.602348+00	\N
af3f64b8-716c-4a1a-85da-d67543d3e588	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	1800	food	2025-11-08 00:00:00+00	Bua Bill	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-09 16:57:42.33479+00	2025-11-09 17:52:57.21584+00	\N
35ca7248-3189-487d-bf4c-e5157677d040	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	\N	expense	2010	other_expense	2025-11-10 00:00:00+00	Tax consultant fee	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-10 08:13:03.209349+00	2025-11-10 08:13:03.209349+00	\N
d7962291-846d-403e-8843-261e01ff1803	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	130	transport	2025-11-10 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-10 08:13:25.444564+00	2025-11-10 08:13:25.444564+00	\N
b1f73683-3341-4a4d-bf91-cd7edfbf6f4d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	140	transport	2025-11-10 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-10 16:00:41.307702+00	2025-11-10 16:00:41.307702+00	\N
5bfbe17e-66ce-4513-8cf6-6335bb8c3004	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	212	donation	2025-11-10 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-10 16:01:27.117587+00	2025-11-10 16:01:27.117587+00	\N
e7573260-58bc-4266-88d6-18bed9968ea2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	320	housing	2025-11-10 00:00:00+00	Water Kettle Repair	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-10 17:50:38.202631+00	2025-11-10 17:50:38.202631+00	\N
2b46ce7c-7c49-4f26-a38f-5a6201c0d81a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	598	utilities	2025-11-10 00:00:00+00	Mobile recharge for Maa	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-10 17:50:10.259861+00	2025-11-10 17:50:10.259861+00	2025-11-10 17:56:24.172274+00
19b78f5a-f034-44aa-b7e4-485ffd8afda5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	75	food	2025-11-07 00:00:00+00	Singara from Bohera	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-10 17:57:25.315287+00	2025-11-10 17:57:25.315287+00	\N
d9045428-a73b-43a1-9b04-538bfa949944	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	140	transport	2025-11-11 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-11 14:12:16.786534+00	2025-11-11 14:12:16.786534+00	\N
b1b86e67-45d7-459c-a136-1c6c7618213c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	130	transport	2025-11-11 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-11 14:12:34.099379+00	2025-11-11 14:12:34.099379+00	\N
502ddf6d-9935-45b9-881d-28bbd91b8b5c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	3625	housing	2025-11-11 00:00:00+00	House rent	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-11 15:55:20.06367+00	2025-11-11 15:55:20.06367+00	\N
dbdb88b3-1c57-4514-b006-a1918b720a29	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	food	2025-11-12 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-12 15:27:36.352366+00	2025-11-12 15:27:36.352366+00	\N
13784bbd-ec16-467e-8d02-a2b5ef9a1a39	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	150	transport	2025-11-12 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-12 15:27:53.485706+00	2025-11-12 15:27:53.485706+00	\N
2675b9d5-850a-477e-b734-e48c4a6e2d4a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-12 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-12 15:28:09.012882+00	2025-11-12 15:28:09.012882+00	\N
a2873c77-6b8f-41ac-9cad-4ac3950169cd	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	125	transport	2025-11-13 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-13 14:42:09.649976+00	2025-11-13 14:42:09.649976+00	\N
df345915-6483-4010-8928-32934cde0601	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	50	food	2025-11-13 00:00:00+00	Bohera breakfast	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-13 14:42:27.943071+00	2025-11-13 14:42:27.943071+00	\N
aa411452-434b-434e-a0ab-ac22ee0c5ea2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-13 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-13 14:44:50.841097+00	2025-11-13 14:44:50.841097+00	\N
178abbaf-f077-475f-ac52-7f6bab746ed0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-13 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-13 14:45:18.775564+00	2025-11-13 14:45:18.775564+00	\N
0f33a456-55b0-4522-8277-7e77dbcbcddb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	10	donation	2025-11-13 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-13 14:46:17.665382+00	2025-11-13 14:46:17.665382+00	\N
e625c60a-636d-4573-9578-dbe2e6315320	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-09 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.109868+00	2025-11-15 08:02:35.109868+00	2025-11-15 15:30:29.779361+00
f2ca6c30-e6c3-4222-9213-2eae63c2769c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	food	2025-11-13 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-13 14:47:01.753096+00	2025-11-13 14:47:01.753096+00	\N
508fbd44-d158-43cf-801e-71bbc89f3cc4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	c3288300-472e-413e-9588-48542b2f66d6	transfer	5000	transfer	2025-11-13 00:00:00+00	Transfer between accounts	[]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-13 17:15:21.952+00	2025-11-13 17:15:21.952+00	\N
d9420559-1b1b-4907-9c55-7b1acbc42833	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	income	800	other_income	2025-11-14 00:00:00+00	Transfer from Rocket to bKash via NPSB	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-14 16:03:12.681769+00	2025-11-14 16:03:12.681769+00	\N
d5c7d302-5ec2-4bfd-ad49-75681d15f6de	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-11 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.112127+00	2025-11-15 08:02:35.112127+00	\N
ee7ba383-3a9a-4f19-8ebc-c5f8cd8927d7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-12 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.116355+00	2025-11-15 08:02:35.116355+00	\N
8f2e3d1f-d984-49fa-b492-635efd195597	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-15 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.119953+00	2025-11-15 08:02:35.119953+00	\N
ef7e519d-2cb1-4e6f-9750-2ac80dacd4dc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	670	food	2025-11-15 00:00:00+00	Tehari Ghor Dhanmondi	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:11:26.480865+00	2025-11-15 15:11:26.480865+00	\N
176c54b1-3bc7-4d2e-9525-51fa1ffe705c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-12 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-12 15:27:21.809844+00	2025-11-13 14:46:34.972559+00	2025-11-15 15:30:09.224551+00
3472a26a-046a-45e9-823e-007907a41983	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-08 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.108081+00	2025-11-15 08:02:35.108081+00	2025-11-15 15:30:52.206478+00
c6245739-4004-4465-b48f-eb0752d2a93b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-10 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.110873+00	2025-11-15 08:02:35.110873+00	2025-11-15 15:31:55.464566+00
2fe39b89-056e-4f53-bfb0-0e63d65ceef9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-14 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.118754+00	2025-11-15 08:02:35.118754+00	2025-11-15 08:03:00.534528+00
4d63308a-ea1c-465e-9a35-eb5d8d47f60f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	135	shopping	2025-11-15 00:00:00+00	From Daraz	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:22:58.596185+00	2025-11-15 15:22:58.596185+00	\N
3f6106fa-36fc-491e-91c5-4a9248f3a334	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	65	food	2025-11-15 00:00:00+00	From aarong	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:23:17.637399+00	2025-11-15 15:23:17.637399+00	\N
ca5c220a-4d8e-4618-a639-4e11c1202aa4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	160	transport	2025-11-15 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:11:58.443156+00	2025-11-15 15:11:58.443156+00	\N
bc2d45c1-fe27-48e8-abae-5b53c9b929bf	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	190	transport	2025-11-15 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:12:23.418446+00	2025-11-15 15:12:23.418446+00	\N
e6e57da4-2aa6-4b5b-9e7f-04c1bf0180a2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	280	transport	2025-11-15 00:00:00+00	CNG fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:13:14.598187+00	2025-11-15 15:13:14.598187+00	\N
7dfa942c-2897-42d8-83d9-ee3ddb69bac4	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	429	shopping	2025-11-15 00:00:00+00	From foodpanda	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:15:03.004998+00	2025-11-15 15:15:03.004998+00	\N
372607e8-c2ba-4620-a9b7-6e0fecb035a8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	654.39	food	2025-11-15 00:00:00+00	From foodpanda	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-15 15:20:30.437874+00	2025-11-15 15:20:30.437874+00	\N
a27ef7d5-3878-42ab-a1ba-e004363051e7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-13 00:00:00+00	Donation in the path of Allah	\N	\N	\N	\N	8845435a-f95f-448b-850a-dfb53371cf08	\N	\N	f	\N	2025-11-15 08:02:35.117366+00	2025-11-15 08:02:35.117366+00	2025-11-15 15:30:02.756362+00
7e82b970-08b0-4ef1-be17-87497b36a457	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-11 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-11 14:11:57.286529+00	2025-11-11 14:11:57.286529+00	2025-11-15 15:30:13.74497+00
ce9b0658-afba-4435-90be-318cece06e26	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	110	transport	2025-11-16 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-16 17:09:55.80935+00	2025-11-16 17:09:55.80935+00	\N
1e78f3b2-8f2e-474d-941f-4a3573a0e3f3	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	50	transport	2025-11-16 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-16 17:10:12.538226+00	2025-11-16 17:10:12.538226+00	\N
e5b30afb-5004-48c1-bd52-7318a3dbea5f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	20	food	2025-11-16 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-16 17:10:26.863838+00	2025-11-16 17:10:26.863838+00	\N
a0464f2d-ac2d-4d6e-b449-8c8d66e66a81	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	130	food	2025-11-16 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-16 17:10:38.879782+00	2025-11-16 17:10:38.879782+00	\N
e47dca06-1702-4b5b-9983-0204f10f69b7	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	income	30000	other_income	2025-11-16 00:00:00+00	From Vivasoft for KPI	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-16 17:11:00.893783+00	2025-11-16 17:11:00.893783+00	\N
4f614975-debf-4fbf-b1d7-b273280c7f79	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-16 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-16 17:11:35.561241+00	2025-11-16 17:11:35.561241+00	\N
f3778271-821f-43bd-9563-acf5cbc01ab0	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-17 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-17 15:00:46.357938+00	2025-11-17 15:00:46.357938+00	\N
5c60f9b8-7a6d-4a17-9a0a-f5b8613b3b1c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	120	transport	2025-11-17 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-17 15:00:29.701568+00	2025-11-17 15:00:52.893738+00	\N
23ea0999-f751-4b4c-bd42-5ac0f2ffbdac	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	500	donation	2025-11-17 00:00:00+00	Gift for humaun bhai	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-17 15:01:38.479116+00	2025-11-17 15:01:38.479116+00	\N
7fbfc35d-0a29-4f65-bc2b-9aa613c5da05	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	30	utilities	2025-11-17 00:00:00+00	bKash Staff for go to ST	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-17 15:02:03.019568+00	2025-11-17 15:02:03.019568+00	\N
26940eb4-ae0b-44b1-9834-2fb7954da753	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	income	500	other_income	2025-11-17 00:00:00+00	bKash HR for night activity	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-17 15:02:31.072028+00	2025-11-17 15:02:31.072028+00	\N
07271880-ccf5-402c-b31f-73337e4553cc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-17 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-17 15:03:00.407013+00	2025-11-17 15:03:00.407013+00	\N
e7f9fd22-5639-448a-9582-47c1132334eb	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	252	food	2025-11-18 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-18 15:48:19.899488+00	2025-11-18 15:48:19.899488+00	\N
c80eface-97e5-4f39-aedc-c5ac489f148c	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	160	transport	2025-11-18 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-18 15:48:52.105019+00	2025-11-18 15:48:52.105019+00	\N
6e521053-d2cd-41c8-8b3c-a3ee93188d3b	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-18 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-18 15:48:33.902504+00	2025-11-18 15:48:59.50279+00	\N
d3860014-cd06-486d-98fc-96bd85a5b809	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-18 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-18 15:49:30.470322+00	2025-11-18 15:49:30.470322+00	\N
584ce5de-1559-4e95-94a9-1d2d7a6ea7a2	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	expense	5750	credit_card_payment	2025-11-17 00:00:00+00	Credit card yearly fee	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-18 15:50:00.22864+00	2025-11-18 15:50:00.22864+00	\N
fa02af77-895d-480a-8993-2fa2e59c71c6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	c3288300-472e-413e-9588-48542b2f66d6	transfer	1200	transfer	2025-11-18 00:00:00+00	Transfer between accounts	[]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-18 17:04:51.299+00	2025-11-18 17:04:51.299+00	\N
72d54b38-9878-414d-967e-14eb55d1578d	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	income	30	other_income	2025-11-18 00:00:00+00	Cashback	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-18 17:05:14.336983+00	2025-11-18 17:05:14.336983+00	\N
7afa2bca-5952-4ad4-99b0-056a7e620803	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	1353	healthcare	2025-11-18 00:00:00+00	Medicine from Aroggo	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-18 17:05:43.05921+00	2025-11-18 17:05:43.05921+00	\N
f89aa192-747d-4cba-8176-8b3758618a46	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-20 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:55:19.682951+00	2025-11-20 14:55:19.682951+00	\N
804ed2b9-d8ae-4d2c-b1ef-6e02a0909878	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	110	transport	2025-11-20 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:55:40.736258+00	2025-11-20 14:55:40.736258+00	\N
1901245e-e128-4cca-9f89-034b42a94946	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	150	transport	2025-11-20 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:55:53.388074+00	2025-11-20 14:55:53.388074+00	\N
ec2a230a-ca06-49af-974a-529a62470d57	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-19 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:56:18.256548+00	2025-11-20 14:56:18.256548+00	\N
477b51d9-636f-4b7c-9ac0-605d0a40d81e	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	transport	2025-11-19 00:00:00+00	Rickshaw fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:56:35.557427+00	2025-11-20 14:56:35.557427+00	\N
3b255be4-a638-4a9c-bcb0-de1c8b41ff77	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	135	food	2025-11-19 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:56:53.558525+00	2025-11-20 14:56:53.558525+00	\N
3c4eee26-56d3-42ce-8a41-338d8357f42f	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	25	food	2025-11-19 00:00:00+00	Bohera	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:58:27.955255+00	2025-11-20 14:58:27.955255+00	\N
fee68dd4-f949-4548-a6dd-816e343096af	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	c3288300-472e-413e-9588-48542b2f66d6	transfer	2000	transfer	2025-11-20 00:00:00+00	Transfer between accounts	[]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-20 14:58:39.845+00	2025-11-20 14:58:39.845+00	\N
a9b1999d-43f1-4ea1-8c23-d28aee9f1bed	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-20 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:59:21.260215+00	2025-11-20 14:59:21.260215+00	\N
622d1295-9608-49e1-ba8c-2d3ad25605c5	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-21 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-22 16:49:33.209628+00	2025-11-22 16:49:33.209628+00	\N
dbf29958-e708-445b-941b-b855c975b2f6	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-22 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-22 16:49:46.361933+00	2025-11-22 16:49:46.361933+00	\N
a60047ea-9997-4169-bb55-0cef6e62f636	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	500	donation	2025-11-22 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-22 16:49:57.753345+00	2025-11-22 16:49:57.753345+00	\N
90a92532-7788-42c8-a978-08783d158143	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	c3288300-472e-413e-9588-48542b2f66d6	transfer	2000	transfer	2025-11-22 00:00:00+00	Transfer between accounts	[]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-22 16:50:19.134+00	2025-11-22 16:50:19.134+00	\N
11cca6b5-f1f2-4913-890f-ce6cf177815a	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	2500	lend	2025-11-21 00:00:00+00	Lent to Nurul Huda	\N	\N	\N	\N	\N	\N	\N	f	\N	2025-11-22 16:51:06.153605+00	2025-11-22 16:51:06.153605+00	\N
397fab37-18fd-4982-8b10-ed97a3956141	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	200	transport	2025-11-22 00:00:00+00	Bike fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-22 16:51:33.22353+00	2025-11-22 16:51:33.22353+00	\N
8ab6b900-bce6-4604-a9f3-4ff4d4eb74e9	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	320	food	2025-11-22 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-22 16:51:49.228094+00	2025-11-22 16:51:49.228094+00	\N
bb332524-84ff-4995-82fd-5abd92866552	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	10	transport	2025-11-22 00:00:00+00	Bus fare	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-22 16:52:04.701721+00	2025-11-22 16:52:04.701721+00	\N
851dc38a-6c3d-4f55-96ae-b0301a9cc6dc	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	112	donation	2025-11-23 00:00:00+00	Donation in the path of Allah	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-23 16:19:54.633506+00	2025-11-23 16:19:54.633506+00	\N
fa731082-60f7-4fa3-9615-e7cf107f9e01	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	c3288300-472e-413e-9588-48542b2f66d6	\N	expense	180	food	2025-11-23 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-23 16:20:33.317104+00	2025-11-23 16:20:33.317104+00	\N
483a1bb7-aba0-4486-b176-5aaba95b5327	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	income	564	other_income	2025-11-20 00:00:00+00	Imsurence	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-23 16:22:05.083573+00	2025-11-23 16:22:05.083573+00	\N
aa1bf759-659b-4379-b7fc-190b5748e083	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	6079ac1c-7a40-45b0-913d-4e2dce11dab6	c3288300-472e-413e-9588-48542b2f66d6	transfer	5000	transfer	2025-11-23 00:00:00+00	Transfer between accounts	[]	\N	\N	\N	\N	\N	\N	f	\N	2025-11-23 16:27:32.552+00	2025-11-23 16:27:32.552+00	\N
f4b9815e-51ad-4ff8-99ad-3900c05626a8	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	cfee8a15-8596-439c-9ef3-53bf5252b9f5	\N	expense	100	food	2025-11-20 00:00:00+00	Snacks	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-20 14:57:07.867418+00	2025-11-23 18:10:31.988729+00	\N
cef0a836-ff4f-46ef-acfe-a43212d253ba	f7ffd4e5-ee01-45b6-b953-4af5541b5c96	338a8a7a-841a-44ca-94c9-1e117406ac70	\N	income	2416	other_income	2025-11-23 00:00:00+00	Insurance	[]	\N	\N	\N	\N	\N	[]	f	\N	2025-11-23 16:22:43.185692+00	2025-11-23 18:20:06.869359+00	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: daybook_user
--

COPY public.users (id, username, email, password, full_name, role, last_login, created_at, updated_at, deleted_at) FROM stdin;
f7ffd4e5-ee01-45b6-b953-4af5541b5c96	shafikshaon	shafikshaon@gmail.com	$2a$10$EQAyfOuKvkQz.huxzO9Zx.zwjVi.VzrFNaew/wWTPlw.iz1r6xOwe	Mohd. Shafikur Rahman	user	2025-11-29 05:30:58.904478+00	2025-11-01 11:57:04.984391+00	2025-11-29 05:30:58.904667+00	\N
\.


--
-- Name: account_types account_types_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.account_types
    ADD CONSTRAINT account_types_pkey PRIMARY KEY (id);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: activity_logs activity_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.activity_logs
    ADD CONSTRAINT activity_logs_pkey PRIMARY KEY (id);


--
-- Name: asset_attachments asset_attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.asset_attachments
    ADD CONSTRAINT asset_attachments_pkey PRIMARY KEY (id);


--
-- Name: assets assets_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_pkey PRIMARY KEY (id);


--
-- Name: bill_payments bill_payments_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.bill_payments
    ADD CONSTRAINT bill_payments_pkey PRIMARY KEY (id);


--
-- Name: bills bills_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.bills
    ADD CONSTRAINT bills_pkey PRIMARY KEY (id);


--
-- Name: budgets budgets_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.budgets
    ADD CONSTRAINT budgets_pkey PRIMARY KEY (id);


--
-- Name: credit_card_payments credit_card_payments_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.credit_card_payments
    ADD CONSTRAINT credit_card_payments_pkey PRIMARY KEY (id);


--
-- Name: credit_card_transactions credit_card_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.credit_card_transactions
    ADD CONSTRAINT credit_card_transactions_pkey PRIMARY KEY (id);


--
-- Name: credit_cards credit_cards_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.credit_cards
    ADD CONSTRAINT credit_cards_pkey PRIMARY KEY (id);


--
-- Name: debt_payments debt_payments_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.debt_payments
    ADD CONSTRAINT debt_payments_pkey PRIMARY KEY (id);


--
-- Name: debt_records debt_records_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.debt_records
    ADD CONSTRAINT debt_records_pkey PRIMARY KEY (id);


--
-- Name: goal_contributions goal_contributions_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.goal_contributions
    ADD CONSTRAINT goal_contributions_pkey PRIMARY KEY (id);


--
-- Name: goal_holdings goal_holdings_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.goal_holdings
    ADD CONSTRAINT goal_holdings_pkey PRIMARY KEY (id);


--
-- Name: goals goals_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.goals
    ADD CONSTRAINT goals_pkey PRIMARY KEY (id);


--
-- Name: lend_payments lend_payments_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.lend_payments
    ADD CONSTRAINT lend_payments_pkey PRIMARY KEY (id);


--
-- Name: lend_records lend_records_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.lend_records
    ADD CONSTRAINT lend_records_pkey PRIMARY KEY (id);


--
-- Name: reconciliation_transactions reconciliation_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.reconciliation_transactions
    ADD CONSTRAINT reconciliation_transactions_pkey PRIMARY KEY (id);


--
-- Name: reconciliations reconciliations_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.reconciliations
    ADD CONSTRAINT reconciliations_pkey PRIMARY KEY (id);


--
-- Name: recurring_transactions recurring_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.recurring_transactions
    ADD CONSTRAINT recurring_transactions_pkey PRIMARY KEY (id, template_id);


--
-- Name: rewards rewards_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.rewards
    ADD CONSTRAINT rewards_pkey PRIMARY KEY (id);


--
-- Name: service_records service_records_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.service_records
    ADD CONSTRAINT service_records_pkey PRIMARY KEY (id);


--
-- Name: settings settings_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (id);


--
-- Name: statements statements_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.statements
    ADD CONSTRAINT statements_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_account_types_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_account_types_deleted_at ON public.account_types USING btree (deleted_at);


--
-- Name: idx_account_types_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_account_types_user_id ON public.account_types USING btree (user_id);


--
-- Name: idx_accounts_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_accounts_deleted_at ON public.accounts USING btree (deleted_at);


--
-- Name: idx_accounts_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_accounts_user_id ON public.accounts USING btree (user_id);


--
-- Name: idx_activity_logs_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_activity_logs_deleted_at ON public.activity_logs USING btree (deleted_at);


--
-- Name: idx_activity_logs_module; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_activity_logs_module ON public.activity_logs USING btree (module);


--
-- Name: idx_activity_logs_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_activity_logs_user_id ON public.activity_logs USING btree (user_id);


--
-- Name: idx_asset_attachments_asset_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_asset_attachments_asset_id ON public.asset_attachments USING btree (asset_id);


--
-- Name: idx_asset_attachments_attachment_type; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_asset_attachments_attachment_type ON public.asset_attachments USING btree (attachment_type);


--
-- Name: idx_asset_attachments_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_asset_attachments_deleted_at ON public.asset_attachments USING btree (deleted_at);


--
-- Name: idx_asset_attachments_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_asset_attachments_user_id ON public.asset_attachments USING btree (user_id);


--
-- Name: idx_assets_category; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_assets_category ON public.assets USING btree (category);


--
-- Name: idx_assets_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_assets_deleted_at ON public.assets USING btree (deleted_at);


--
-- Name: idx_assets_purchase_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_assets_purchase_date ON public.assets USING btree (purchase_date);


--
-- Name: idx_assets_status; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_assets_status ON public.assets USING btree (status);


--
-- Name: idx_assets_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_assets_user_id ON public.assets USING btree (user_id);


--
-- Name: idx_bill_payments_bill_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_bill_payments_bill_id ON public.bill_payments USING btree (bill_id);


--
-- Name: idx_bill_payments_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_bill_payments_deleted_at ON public.bill_payments USING btree (deleted_at);


--
-- Name: idx_bill_payments_payment_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_bill_payments_payment_date ON public.bill_payments USING btree (payment_date);


--
-- Name: idx_bill_payments_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_bill_payments_user_id ON public.bill_payments USING btree (user_id);


--
-- Name: idx_bills_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_bills_deleted_at ON public.bills USING btree (deleted_at);


--
-- Name: idx_bills_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_bills_user_id ON public.bills USING btree (user_id);


--
-- Name: idx_budgets_category_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_budgets_category_id ON public.budgets USING btree (category_id);


--
-- Name: idx_budgets_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_budgets_deleted_at ON public.budgets USING btree (deleted_at);


--
-- Name: idx_budgets_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_budgets_user_id ON public.budgets USING btree (user_id);


--
-- Name: idx_credit_card_payments_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_payments_account_id ON public.credit_card_payments USING btree (account_id);


--
-- Name: idx_credit_card_payments_card_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_payments_card_id ON public.credit_card_payments USING btree (card_id);


--
-- Name: idx_credit_card_payments_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_payments_deleted_at ON public.credit_card_payments USING btree (deleted_at);


--
-- Name: idx_credit_card_payments_transaction_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_payments_transaction_id ON public.credit_card_payments USING btree (transaction_id);


--
-- Name: idx_credit_card_payments_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_payments_user_id ON public.credit_card_payments USING btree (user_id);


--
-- Name: idx_credit_card_transactions_card_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_transactions_card_id ON public.credit_card_transactions USING btree (card_id);


--
-- Name: idx_credit_card_transactions_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_transactions_deleted_at ON public.credit_card_transactions USING btree (deleted_at);


--
-- Name: idx_credit_card_transactions_transaction_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_transactions_transaction_id ON public.credit_card_transactions USING btree (transaction_id);


--
-- Name: idx_credit_card_transactions_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_card_transactions_user_id ON public.credit_card_transactions USING btree (user_id);


--
-- Name: idx_credit_cards_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_cards_deleted_at ON public.credit_cards USING btree (deleted_at);


--
-- Name: idx_credit_cards_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_credit_cards_user_id ON public.credit_cards USING btree (user_id);


--
-- Name: idx_debt_payments_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_payments_account_id ON public.debt_payments USING btree (account_id);


--
-- Name: idx_debt_payments_debt_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_payments_debt_id ON public.debt_payments USING btree (debt_id);


--
-- Name: idx_debt_payments_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_payments_deleted_at ON public.debt_payments USING btree (deleted_at);


--
-- Name: idx_debt_payments_payment_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_payments_payment_date ON public.debt_payments USING btree (payment_date);


--
-- Name: idx_debt_payments_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_payments_user_id ON public.debt_payments USING btree (user_id);


--
-- Name: idx_debt_records_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_records_account_id ON public.debt_records USING btree (account_id);


--
-- Name: idx_debt_records_borrowed_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_records_borrowed_date ON public.debt_records USING btree (borrowed_date);


--
-- Name: idx_debt_records_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_records_deleted_at ON public.debt_records USING btree (deleted_at);


--
-- Name: idx_debt_records_status; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_records_status ON public.debt_records USING btree (status);


--
-- Name: idx_debt_records_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_debt_records_user_id ON public.debt_records USING btree (user_id);


--
-- Name: idx_goal_contributions_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_contributions_date ON public.goal_contributions USING btree (date);


--
-- Name: idx_goal_contributions_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_contributions_deleted_at ON public.goal_contributions USING btree (deleted_at);


--
-- Name: idx_goal_contributions_goal_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_contributions_goal_id ON public.goal_contributions USING btree (goal_id);


--
-- Name: idx_goal_contributions_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_contributions_user_id ON public.goal_contributions USING btree (user_id);


--
-- Name: idx_goal_holdings_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_holdings_deleted_at ON public.goal_holdings USING btree (deleted_at);


--
-- Name: idx_goal_holdings_goal_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_holdings_goal_id ON public.goal_holdings USING btree (goal_id);


--
-- Name: idx_goal_holdings_type; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_holdings_type ON public.goal_holdings USING btree (type);


--
-- Name: idx_goal_holdings_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goal_holdings_user_id ON public.goal_holdings USING btree (user_id);


--
-- Name: idx_goals_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goals_deleted_at ON public.goals USING btree (deleted_at);


--
-- Name: idx_goals_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_goals_user_id ON public.goals USING btree (user_id);


--
-- Name: idx_lend_payments_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_payments_account_id ON public.lend_payments USING btree (account_id);


--
-- Name: idx_lend_payments_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_payments_deleted_at ON public.lend_payments USING btree (deleted_at);


--
-- Name: idx_lend_payments_lend_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_payments_lend_id ON public.lend_payments USING btree (lend_id);


--
-- Name: idx_lend_payments_payment_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_payments_payment_date ON public.lend_payments USING btree (payment_date);


--
-- Name: idx_lend_payments_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_payments_user_id ON public.lend_payments USING btree (user_id);


--
-- Name: idx_lend_records_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_records_account_id ON public.lend_records USING btree (account_id);


--
-- Name: idx_lend_records_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_records_deleted_at ON public.lend_records USING btree (deleted_at);


--
-- Name: idx_lend_records_lent_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_records_lent_date ON public.lend_records USING btree (lent_date);


--
-- Name: idx_lend_records_status; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_records_status ON public.lend_records USING btree (status);


--
-- Name: idx_lend_records_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_lend_records_user_id ON public.lend_records USING btree (user_id);


--
-- Name: idx_reconciliation_transactions_reconciliation_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_reconciliation_transactions_reconciliation_id ON public.reconciliation_transactions USING btree (reconciliation_id);


--
-- Name: idx_reconciliation_transactions_transaction_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_reconciliation_transactions_transaction_id ON public.reconciliation_transactions USING btree (transaction_id);


--
-- Name: idx_reconciliations_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_reconciliations_account_id ON public.reconciliations USING btree (account_id);


--
-- Name: idx_reconciliations_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_reconciliations_deleted_at ON public.reconciliations USING btree (deleted_at);


--
-- Name: idx_reconciliations_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_reconciliations_user_id ON public.reconciliations USING btree (user_id);


--
-- Name: idx_recurring_transactions_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_account_id ON public.recurring_transactions USING btree (template_account_id);


--
-- Name: idx_recurring_transactions_category_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_category_id ON public.recurring_transactions USING btree (template_category_id);


--
-- Name: idx_recurring_transactions_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_date ON public.recurring_transactions USING btree (template_date);


--
-- Name: idx_recurring_transactions_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_deleted_at ON public.recurring_transactions USING btree (template_deleted_at, deleted_at);


--
-- Name: idx_recurring_transactions_fixed_deposit_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_fixed_deposit_id ON public.recurring_transactions USING btree (template_fixed_deposit_id);


--
-- Name: idx_recurring_transactions_investment_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_investment_id ON public.recurring_transactions USING btree (template_investment_id);


--
-- Name: idx_recurring_transactions_reconciled; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_reconciled ON public.recurring_transactions USING btree (template_reconciled);


--
-- Name: idx_recurring_transactions_savings_goal_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_savings_goal_id ON public.recurring_transactions USING btree (template_savings_goal_id);


--
-- Name: idx_recurring_transactions_to_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_to_account_id ON public.recurring_transactions USING btree (template_to_account_id);


--
-- Name: idx_recurring_transactions_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_recurring_transactions_user_id ON public.recurring_transactions USING btree (user_id, template_user_id);


--
-- Name: idx_rewards_card_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_rewards_card_id ON public.rewards USING btree (card_id);


--
-- Name: idx_rewards_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_rewards_deleted_at ON public.rewards USING btree (deleted_at);


--
-- Name: idx_rewards_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_rewards_user_id ON public.rewards USING btree (user_id);


--
-- Name: idx_service_records_asset_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_service_records_asset_id ON public.service_records USING btree (asset_id);


--
-- Name: idx_service_records_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_service_records_deleted_at ON public.service_records USING btree (deleted_at);


--
-- Name: idx_service_records_service_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_service_records_service_date ON public.service_records USING btree (service_date);


--
-- Name: idx_service_records_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_service_records_user_id ON public.service_records USING btree (user_id);


--
-- Name: idx_settings_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_settings_deleted_at ON public.settings USING btree (deleted_at);


--
-- Name: idx_settings_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE UNIQUE INDEX idx_settings_user_id ON public.settings USING btree (user_id);


--
-- Name: idx_statements_card_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_statements_card_id ON public.statements USING btree (card_id);


--
-- Name: idx_statements_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_statements_deleted_at ON public.statements USING btree (deleted_at);


--
-- Name: idx_statements_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_statements_user_id ON public.statements USING btree (user_id);


--
-- Name: idx_tags_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_tags_deleted_at ON public.tags USING btree (deleted_at);


--
-- Name: idx_tags_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_tags_user_id ON public.tags USING btree (user_id);


--
-- Name: idx_transactions_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_account_id ON public.transactions USING btree (account_id);


--
-- Name: idx_transactions_category_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_category_id ON public.transactions USING btree (category_id);


--
-- Name: idx_transactions_date; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_date ON public.transactions USING btree (date);


--
-- Name: idx_transactions_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_deleted_at ON public.transactions USING btree (deleted_at);


--
-- Name: idx_transactions_fixed_deposit_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_fixed_deposit_id ON public.transactions USING btree (fixed_deposit_id);


--
-- Name: idx_transactions_investment_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_investment_id ON public.transactions USING btree (investment_id);


--
-- Name: idx_transactions_reconciled; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_reconciled ON public.transactions USING btree (reconciled);


--
-- Name: idx_transactions_savings_goal_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_savings_goal_id ON public.transactions USING btree (savings_goal_id);


--
-- Name: idx_transactions_to_account_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_to_account_id ON public.transactions USING btree (to_account_id);


--
-- Name: idx_transactions_user_id; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_transactions_user_id ON public.transactions USING btree (user_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: daybook_user
--

CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: activity_logs fk_activity_logs_user; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.activity_logs
    ADD CONSTRAINT fk_activity_logs_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: asset_attachments fk_assets_attachments; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.asset_attachments
    ADD CONSTRAINT fk_assets_attachments FOREIGN KEY (asset_id) REFERENCES public.assets(id);


--
-- Name: service_records fk_assets_service_records; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.service_records
    ADD CONSTRAINT fk_assets_service_records FOREIGN KEY (asset_id) REFERENCES public.assets(id);


--
-- Name: goal_contributions fk_goals_contributions; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.goal_contributions
    ADD CONSTRAINT fk_goals_contributions FOREIGN KEY (goal_id) REFERENCES public.goals(id);


--
-- Name: goal_holdings fk_goals_holdings; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.goal_holdings
    ADD CONSTRAINT fk_goals_holdings FOREIGN KEY (goal_id) REFERENCES public.goals(id);


--
-- Name: reconciliation_transactions fk_reconciliation_transactions_transaction; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.reconciliation_transactions
    ADD CONSTRAINT fk_reconciliation_transactions_transaction FOREIGN KEY (transaction_id) REFERENCES public.transactions(id);


--
-- Name: reconciliations fk_reconciliations_account; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.reconciliations
    ADD CONSTRAINT fk_reconciliations_account FOREIGN KEY (account_id) REFERENCES public.accounts(id);


--
-- Name: reconciliation_transactions fk_reconciliations_transactions; Type: FK CONSTRAINT; Schema: public; Owner: daybook_user
--

ALTER TABLE ONLY public.reconciliation_transactions
    ADD CONSTRAINT fk_reconciliations_transactions FOREIGN KEY (reconciliation_id) REFERENCES public.reconciliations(id);


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO daybook_user;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO daybook_user;


--
-- PostgreSQL database dump complete
--

\unrestrict gA1daJdzDra9ovXP0herlMbqmJWuZ7zPO9h9FI1CNztDdQIVUHnKH9zCcOASz3w

