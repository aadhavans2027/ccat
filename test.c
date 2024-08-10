//go:build exclude
#include "easysock.h"


int create_socket(int network, char transport) {
	int domain;
	int type;

	if (network == 4) {
		domain = AF_INET;
	} else if (network == 6) {
		domain = AF_INET6;
	} else {
		return -1;
	}

	if (transport == 'T') {
		type = SOCK_STREAM;
	} else if (transport == 'U') {
		type = SOCK_DGRAM;
	} else {
		return -1;
	}

	int newSock = socket(domain,type,0);
	return newSock;
}


int create_addr(int network, char* address, int port,struct sockaddr* dest) {
	if (network == 4) {
		struct sockaddr_in listen_address;

		listen_address.sin_family = AF_INET;
		listen_address.sin_port = htons(port);
		inet_pton(AF_INET,address,&listen_address.sin_addr);
		memcpy(dest,&listen_address,sizeof(listen_address));
		return 0;

	} else if (network == 6) {
		struct sockaddr_in6 listen_ipv6;
		listen_ipv6.sin6_family = AF_INET6;
		listen_ipv6.sin6_port = htons(port);
		inet_pton(AF_INET6,address,&listen_ipv6.sin6_addr);
		memcpy(dest,&listen_ipv6,sizeof(listen_ipv6));
		return 0;

	} else {
		return -202;
	}



}

int create_local (int network, char transport, char* address, int port,struct sockaddr* addr_struct) {
	int socket = create_socket(network,transport);
	if (socket < 0) {
		return (-1 * errno);
	}
	create_addr(network,address,port,addr_struct);
	int addrlen;
	if (network == 4) {
		addrlen = sizeof(struct sockaddr_in);
	} else if (network == 6) {
		addrlen = sizeof(struct sockaddr_in6);
	} else {
		return -202;
	}

	/* The value of addrlen should be the size of the 'sockaddr'.
	This should be set to the size of 'sockaddr_in' for IPv4, and 'sockaddr_in6' for IPv6.
	See https://stackoverflow.com/questions/73707162/socket-bind-failed-with-invalid-argument-error-for-program-running-on-macos */

	int i = bind (socket,addr_struct,(socklen_t)addrlen);
	if (i < 0) {
		return (-1 * errno);
	}
	return socket;
}

int create_remote (int network,char transport,char* address,int port,struct sockaddr* remote_addr_struct) {

	struct addrinfo hints;  /* Used to tell getaddrinfo what kind of address we want */
	struct addrinfo* results; /* Used by getaddrinfo to store the addresses */


	if (check_ip_ver(address) < 0) { /* If the address is a domain name */
		int err_code;
		char* port_str = malloc(10 * sizeof(char));

		sprintf(port_str,"%d",port); /* getaddrinfo expects a string for its port */


		memset(&hints,'\0',sizeof(hints));
		hints.ai_socktype = char_to_socktype(transport);

		err_code = getaddrinfo(address,port_str,&hints,&results);
		if (err_code != 0) {
			return (-1 * err_code);
		}
		remote_addr_struct = results->ai_addr;
		network = inet_to_int(results->ai_family);
	} else {
	        create_addr(network,address,port,remote_addr_struct);
	}

	int socket = create_socket(network,transport);
	if (socket < 0) {
                return (-1 * errno);
        }

       	int addrlen;
	if (network == 4) {
		addrlen = sizeof(struct sockaddr_in);
	} else if (network == 6) {
		addrlen = sizeof(struct sockaddr_in6);
	} else {
		return (-202);
	}

	/* The value of addrlen should be the size of the 'sockaddr'.
	This should be set to the size of 'sockaddr_in' for IPv4, and 'sockaddr_in6' for IPv6.
	See https://stackoverflow.com/questions/73707162/socket-bind-failed-with-invalid-argument-error-for-program-running-on-macos */

        int i = connect(socket,remote_addr_struct,(socklen_t)addrlen);
	if (i < 0) {
		return (-1 * errno);
	}
        return socket;
}


int check_ip_ver(char* address) {
	char buffer[16]; /* 16 chars - 128 bits - is enough to hold an ipv6 address */
	if (inet_pton(AF_INET,address,buffer) == 1) {
		return 4;
	} else if (inet_pton(AF_INET6,address,buffer) == 1) {
		return 6;
	} else {
		return -1;
	}
}

int int_to_inet(int network) {
	if (network == 4) {
		return AF_INET;
	} else if (network == 6) {
		return AF_INET6;
	} else {
		return -202;
	}
}

int inet_to_int(int af_type) {
	if (af_type == AF_INET) {
		return 4;
	} else if (af_type == AF_INET6) {
		return 6;
	} else {
		return -207;
	}
}

int char_to_socktype(char transport) {
	if (transport == 'T') {
		return SOCK_STREAM;
	} else if (transport == 'U') {
		return SOCK_DGRAM;
	} else {
		return -250;
	}
}
