#[cfg(test)]
mod test {
    use core::result::ResultTrait;
    use core::option::OptionTrait;
    use core::traits::TryInto;
    use openzeppelin::{
        utils::serde::SerializedAppend,
        token::erc20::interface::{IERC20Dispatcher, IERC20DispatcherTrait},
        account::interface::{AccountABIDispatcherTrait, AccountABIDispatcher}, tests::utils,
    };
    use starknet::{
        ContractAddress, ClassHash, contract_address_const, get_contract_address,
        testing::{
            set_contract_address, set_caller_address, set_signature, set_transaction_hash,
            set_version
        },
        account::Call,
    };

    const ERC20_TEST_CLASS_HASH: felt252 =
        0xfa15f33d9a964602972ee0635ba5e641646f0944d7dc279360e7ec943dce6a;
    const ACCOUNT_TEST_CLASS_HASH: felt252 =
        0xd5ad229820cc3391b5d3888c6ce1e08f010ce0d5be429e8030dfc603c60dc8;

    fn deploy_nekomoto() {}
    fn deploy_nekocoin() {}
    fn deploy_prism() {}
    fn deploy_shard() {}

    fn deploy_account(salt: felt252) -> AccountABIDispatcher {
        set_version(1);

        let mut calldata = array![];
        set_signature(
            array![
                0x6bc22689efcaeacb9459577138aff9f0af5b77ee7894cdc8efabaf760f6cf6e,
                0x295989881583b9325436851934334faa9d639a2094cd1e2f8691c8a71cd4cdf
            ]
                .span()
        );
        set_transaction_hash(0x601d3d2e265c10ff645e1554c435e72ce6721f0ba5fc96f0c650bfc6231191a);
        calldata.append(0x26da8d11938b76025862be14fdb8b28438827f73e75e86f7bfa38b196951fa7);

        let address = deploy_with_salt(ACCOUNT_TEST_CLASS_HASH, calldata, salt);
        AccountABIDispatcher { contract_address: address }
    }

    fn deploy_with_salt(
        classhash: felt252, calldata: Array<felt252>, salt: felt252
    ) -> ContractAddress {
        let (address, _) = starknet::syscalls::deploy_syscall(
            classhash.try_into().unwrap(), salt, calldata.span(), false
        )
            .expect('deploy failed');
        address
    }

    fn deploy(classhash: felt252, calldata: Array<felt252>) -> ContractAddress {
        let (address, _) = starknet::syscalls::deploy_syscall(
            classhash.try_into().unwrap(), 0, calldata.span(), false
        )
            .expect('deploy failed');
        address
    }
}

