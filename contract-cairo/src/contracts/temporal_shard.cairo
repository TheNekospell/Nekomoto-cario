#[starknet::contract]
pub mod TemporalShard {
    use openzeppelin::token::erc721::erc721::ERC721Component::InternalTrait;
    use openzeppelin::introspection::src5::SRC5Component;
    use openzeppelin::token::erc721::{ERC721Component, ERC721HooksEmptyImpl};
    use starknet::{ContractAddress, get_caller_address};

    component!(path: ERC721Component, storage: erc721, event: ERC721Event);
    component!(path: SRC5Component, storage: src5, event: SRC5Event);

    #[abi(embed_v0)]
    impl ERC721MixinImpl = ERC721Component::ERC721MixinImpl<ContractState>;
    impl ERC721InternalImpl = ERC721Component::InternalImpl<ContractState>;

    #[storage]
    struct Storage {
        host: ContractAddress,
        token_id: u256,
        #[substorage(v0)]
        erc721: ERC721Component::Storage,
        #[substorage(v0)]
        src5: SRC5Component::Storage
    }

    #[event]
    #[derive(Drop, starknet::Event)]
    enum Event {
        #[flat]
        ERC721Event: ERC721Component::Event,
        #[flat]
        SRC5Event: SRC5Component::Event
    }

    #[constructor]
    fn constructor(ref self: ContractState, host: ContractAddress) {
        let name = "TemporalShard";
        let symbol = "TemporalShard";
        let base_uri = "TBD";

        self.token_id.write(1);
        self.host.write(host);

        self.erc721.initializer(name, symbol, base_uri);
    }

    #[external(v0)]
    fn mint(ref self: ContractState, recipient: ContractAddress, count: u256) {
        assert(get_caller_address() == self.host.read(), 'Only the host can mint');
        let mut i = 0;
        loop {
            if i == count {
                break;
            }

            self.erc721.mint(recipient, self.token_id.read());
            self.token_id.write(self.token_id.read() + 1);

            i = i + 1;
        }
    }

    #[external(v0)]
    fn burn(ref self: ContractState, token_id: u256) {
        self.erc721.burn(token_id);
    }

    #[external(v0)]
    fn burnFrom(ref self: ContractState, account: ContractAddress, token_id: u256) {
        self.erc721._check_authorized(account, get_caller_address(), token_id);
        self.erc721.burn(token_id);
    }
}
